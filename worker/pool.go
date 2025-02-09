package worker

import (
	"context"
	"distributed-job-processor/job"
	"distributed-job-processor/queue"
	"fmt"
	"log"
	"sync"
	"time"
)

type WorkerPool struct {
	workers    []*Worker
	maxWorkers int
	jobQueue   queue.Queue
	ctx        context.Context
	cancel     context.CancelFunc
}

// Create a new worker pool
func NewWorkerPool(maxWorkers int, jobQueue queue.Queue) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	if maxWorkers <= 0 {
		maxWorkers = 1 // Ensure at least one worker
	}
	return &WorkerPool{
		workers:    make([]*Worker, 0, maxWorkers), // 0 length but maxWorkers capacity
		maxWorkers: maxWorkers,
		jobQueue:   jobQueue,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (wp *WorkerPool) distributeJobs() {
	log.Println("Starting job distribution")
	for {
		select {
		case <-wp.ctx.Done():
			log.Println("Stopping Job distribution")
			return
		default:
			log.Println("Checking queue for jobs...")
			job, err := wp.jobQueue.DequeueJob()
			if err != nil {
				log.Printf("Error popping job: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}

			log.Printf("Got job %s from queue, looking for idle worker", job.ID)
			workerFound := false
			for _, worker := range wp.workers {
				log.Printf("Checking worker %s - State: %s", worker.ID, worker.state)
				if worker.state == StateIdle {
					worker.state = StateWorking
					log.Printf("Assigning job %s to worker %s", job.ID, worker.ID)
					worker.jobChan <- job
					workerFound = true
					break
				}
			}
			if !workerFound {
				log.Printf("No idle workers, requeueing job %s", job.ID)
				if err := wp.jobQueue.EnqueueJob(job); err != nil {
					log.Printf("Failed to requeue job: %v", err)
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// Start the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.maxWorkers; i++ {
		worker := NewWorker(fmt.Sprintf("worker-%d", i), wp.jobQueue)
		wp.workers = append(wp.workers, worker)
		worker.Start()
	}

	go wp.distributeJobs()
}

func (wp *WorkerPool) Stop() {
	log.Println("Stopping worker pool...")
	wp.cancel() // Signal all workers to stop

	// Wait for all workers to stop
	var wg sync.WaitGroup
	for _, worker := range wp.workers {
		wg.Add(1)
		go func(w *Worker) {
			defer wg.Done()
			w.Stop()
		}(worker)
	}

	// Wait for all workers to finish
	wg.Wait()
	log.Println("All workers stopped")
}

func (wp *WorkerPool) SubmitJob(job *job.Job) error {
	log.Printf("Submitting new job with id %s to queue", job.ID)
	return wp.jobQueue.EnqueueJob(job)
}
