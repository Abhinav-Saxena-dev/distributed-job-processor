package queue

import (
	"distributed-job-processor/job"
	"errors"
	"sync"
)

type Queue interface {
	Length() int
	isQueueJobsEmpty() bool
	Push(job *job.Job) error
	Pop() (*job.Job, error)
	Peek() (*job.Job, error)
}

type queue struct {
	jobs []*job.Job
	mu   sync.Mutex
}

func NewQueue() Queue {
	return &queue{
		jobs: make([]*job.Job, 0),
		// mu is automatically initialized
	}
}

func (queue *queue) Length() int {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	return len(queue.jobs)
}

func (queue *queue) isQueueJobsEmpty() bool {
	return len(queue.jobs) == 0
}

// Push add element to the end of slice
func (queue *queue) Push(job *job.Job) error {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if job == nil {
		return errors.New("cannot push nil job")
	}
	queue.jobs = append(queue.jobs, job)
	return nil
}

// Pop Return first element of the slice
func (queue *queue) Pop() (*job.Job, error) {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if queue.isQueueJobsEmpty() {
		return nil, errors.New("Queue is Empty")
	}
	poppedJob := queue.jobs[0]
	queue.jobs = queue.jobs[1:]
	return poppedJob, nil
}

func (queue *queue) Peek() (*job.Job, error) {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if queue.isQueueJobsEmpty() {
		return nil, errors.New("queue is empty")
	}
	return queue.jobs[0], nil
}
