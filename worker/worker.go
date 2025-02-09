package worker

import (
	"context"
	"distributed-job-processor/job"
	"distributed-job-processor/queue"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

// Read about each package having its own constants even if redundant.

type WorkerState string

const (
	StateIdle    WorkerState = "idle"
	StateWorking WorkerState = "working"
	StateStopped WorkerState = "stopped"
	maxRetries               = 3
	jobTimeout               = 10 * time.Second
)

// Metrics and Progress Tracking
type WorkerMetrics struct {
	JobsProcessed       int64
	JobsFailed          int64
	TotalProcessingTime time.Duration
}

type Worker struct {
	ID      string
	jobChan chan *job.Job
	quit    chan bool
	done    chan bool
	state   WorkerState
	metrics WorkerMetrics
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewWorker(id string, jobQueue queue.Queue) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		ID:      id,
		jobChan: make(chan *job.Job),
		quit:    make(chan bool),
		done:    make(chan bool),
		state:   StateIdle,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start worker in a separate goroutine
func (w *Worker) Start() {
	go w.processJobs()
}

func (w *Worker) Stop() {
	log.Printf("Stopping worker %s...", w.ID)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	w.cancel()

	select {
	case <-w.done:
		log.Printf("Worker %s stopped gracefully", w.ID)
	case <-ctx.Done():
		log.Printf("Worker %s force stopped after timeout", w.ID)
	}
}

func (w *Worker) updateJobProgress(j *job.Job, progress float64) {
	j.JobProgress.Progress = progress
	log.Printf("Job %s progress: %.2f%%", j.ID, progress)
}

func (w *Worker) processJobs() {
	log.Printf("Worker %s started", w.ID)
	for {
		select {
		case <-w.ctx.Done():
			log.Printf("Worker %s context cancelled", w.ID)
			w.state = StateStopped
			w.done <- true
			return
		case processJob := <-w.jobChan:
			log.Printf("Worker %s processing job %s of type %s", w.ID, processJob.ID, processJob.Type)
			w.processJobWithMetrics(processJob)
			w.state = StateIdle
		}
	}
}

func (w *Worker) processJobWithMetrics(processJob *job.Job) {
	if err := processJob.Validate(); err != nil {
		log.Printf("Invalid job: %v", err)
		processJob.JobProgress.Status = job.StatusFailed
		atomic.AddInt64(&w.metrics.JobsFailed, 1)
		return
	}
	startTime := time.Now()
	processJob.JobProgress.StartedAt = startTime
	w.updateJobProgress(processJob, 0)
	if err := w.processJobWithRetries(processJob); err != nil {
		log.Printf("Worker %s failed to process job %s: %v", w.ID, processJob.ID, err)
		processJob.JobProgress.Status = job.StatusFailed
		atomic.AddInt64(&w.metrics.JobsFailed, 1)
	} else {
		log.Printf("Worker %s completed job %s successfully", w.ID, processJob.ID)
		atomic.AddInt64(&w.metrics.JobsProcessed, 1)
	}

	processJob.JobProgress.FinishedAt = time.Now()
	atomic.AddInt64((*int64)(&w.metrics.TotalProcessingTime),
		int64(time.Since(startTime)))
}

func (w *Worker) processJobWithRetries(processJob *job.Job) error {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Worker %s attempting job %s (attempt %d/%d)", w.ID, processJob.ID, attempt, maxRetries)
		w.updateJobProgress(processJob, float64(attempt-1)/(float64(maxRetries))*100)
		err := w.processJobWithTimeout(processJob)
		if err == nil {
			w.updateJobProgress(processJob, 100)
			return nil
		}
		lastErr = err
		log.Printf("Attempt %d failed: %v", attempt, err)
		if attempt < maxRetries {
			backoff := time.Duration(attempt) * time.Second
			time.Sleep(backoff)
		}
	}
	return fmt.Errorf("failed after %d attempts. Last error: %v", maxRetries, lastErr)
}

func (w *Worker) processJobWithTimeout(processJob *job.Job) error {
	ctx, cancel := context.WithTimeout(w.ctx, jobTimeout)
	defer cancel()
	done := make(chan error)
	go func() {
		done <- w.processJob(processJob)
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("job timed out after %v", jobTimeout)
		}
		return ctx.Err()
	}
}

func (w *Worker) processJob(processJob *job.Job) error {
	processJob.JobProgress.Status = job.StatusRunning
	// Process based on job type
	switch processJob.Type {
	case job.EmailJob:
		return w.processEmailJob(processJob)
	case job.APIJob:
		return w.processAPIJob(processJob)
	case job.DataJobType:
		return w.processDataJob(processJob)
	default:
		processJob.JobProgress.Status = job.StatusFailed
		return fmt.Errorf("unsupported job type: %s", processJob.Type)
	}
}

func (w *Worker) processEmailJob(processJob *job.Job) error {
	// Simulate email sending
	time.Sleep(5 * time.Second)
	processJob.JobProgress.Status = job.StatusCompleted
	log.Printf("Worker %s has completed the job with ID: %s", w.ID, processJob.ID)
	return nil
}

func (w *Worker) processAPIJob(processJob *job.Job) error {
	// Simulate API calling
	time.Sleep(5 * time.Second)
	processJob.JobProgress.Status = job.StatusCompleted
	log.Printf("Worker %s has completed the job with ID: %s", w.ID, processJob.ID)
	return nil
}

func (w *Worker) processDataJob(processJob *job.Job) error {
	// Simulate data job processing
	time.Sleep(5 * time.Second)
	processJob.JobProgress.Status = job.StatusCompleted
	log.Printf("Worker %s has completed the job with ID: %s", w.ID, processJob.ID)
	return nil
}

/* --------------------------------------------------- LEGACY CODE --------------------------------------------------- */
//package worker
//
//import (
//"distributed-job-processor/job"
//"distributed-job-processor/queue"
//"fmt"
//"log"
//"time"
//)
//
//type WorkerState string
//
//const (
//	StateIdle    WorkerState = "idle"
//	StateWorking WorkerState = "working"
//	StateStopped WorkerState = "stopped"
//	maxRetries               = 3
//	jobTimeout               = 5 * time.Second
//)
//
//type Worker struct {
//	ID       string
//	jobQueue *queue.Queue
//	quit     chan bool
//	done     chan bool
//	state    WorkerState
//}
//
//func NewWorker(id string, jobQueue *queue.Queue) *Worker {
//	return &Worker{
//		ID:       id,
//		jobQueue: jobQueue,
//		quit:     make(chan bool),
//		done:     make(chan bool),
//		state:    StateIdle,
//	}
//}
//
//// Start worker in a separate goroutine
//func (w *Worker) Start() {
//	go w.processJobs()
//}
//
//func (w *Worker) Stop() {
//	w.quit <- true
//	<-w.done
//	log.Printf("Worker %s stopped", w.ID)
//}
//
//func (w *Worker) processJobs() {
//	log.Printf("Worker %s started", w.ID)
//	for {
//		select {
//		case <-w.quit:
//			log.Printf("Worker %s shutting down", w.ID)
//			w.state = StateStopped
//			w.done <- true
//			return
//		default:
//			w.state = StateIdle
//			processJob, err := w.jobQueue.DequeueJob()
//			if err != nil {
//				time.Sleep(time.Second)
//				continue
//			}
//			w.state = StateWorking
//			log.Printf("Worker %s processing job %s of type %s", w.ID, processJob.ID, processJob.Type)
//			if err := w.processJobWithRetries(processJob); err != nil {
//				log.Printf("Worker %s failed to process job %s: %v", w.ID, processJob.ID, err)
//				processJob.Status = job.StatusFailed
//			} else {
//				log.Printf("Worker %s completed job %s successfully", w.ID, processJob.ID)
//			}
//		}
//	}
//}
//
//func (w *Worker) processJobWithRetries(processJob *job.Job) error {
//	var lastErr error
//	for attempt := 1; attempt <= maxRetries; attempt++ {
//		log.Printf("Worker %s attempting job %s (attempt %d/%d)", w.ID, processJob.ID, attempt, maxRetries)
//		err := w.processJobWithTimeout(processJob)
//		if err == nil {
//			return nil
//		}
//		lastErr = err
//		log.Printf("Attempt %d failed: %v", attempt, err)
//		if attempt < maxRetries {
//			backoff := time.Duration(attempt) * time.Second
//			time.Sleep(backoff)
//		}
//	}
//	return fmt.Errorf("failed after %d attempts. Last error: %v", maxRetries, lastErr)
//}
//
//func (w *Worker) processJobWithTimeout(processJob *job.Job) error {
//	done := make(chan error)
//	go func() {
//		done <- w.processJob(processJob)
//	}()
//	select {
//	case err := <-done:
//		return err
//	case <-time.After(jobTimeout): // If timeout happens first
//		// Go runtime automatically sends a value after jobTimeout duration
//		// You don't send this value - Go does it for you
//		return fmt.Errorf("job timed out after %v", jobTimeout)
//	}
//}
//
//func (w *Worker) processJob(processJob *job.Job) error {
//	processJob.Status = job.StatusRunning
//	// Process based on job type
//	switch processJob.Type {
//	case job.EmailJob:
//		return w.processEmailJob(processJob)
//	case job.APIJob:
//		return w.processAPIJob(processJob)
//	case job.DataJobType:
//		return w.processDataJob(processJob)
//	default:
//		processJob.Status = job.StatusFailed
//		return fmt.Errorf("unsupported job type: %s", processJob.Type)
//	}
//}
//
//func (w *Worker) processEmailJob(processJob *job.Job) error {
//	// Simulate email sending
//	time.Sleep(2 * time.Second)
//	processJob.Status = job.StatusCompleted
//	return nil
//}
//
//func (w *Worker) processAPIJob(processJob *job.Job) error {
//	// Simulate API calling
//	time.Sleep(2 * time.Second)
//	processJob.Status = job.StatusCompleted
//	return nil
//}
//
//func (w *Worker) processDataJob(processJob *job.Job) error {
//	// Simulate data job processing
//	time.Sleep(2 * time.Second)
//	processJob.Status = job.StatusCompleted
//	return nil
//}
