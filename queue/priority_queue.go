package queue

import (
	"container/heap"
	"distributed-job-processor/job"
	"errors"
	"log"
	"sync"
)

type PriorityQueue struct {
	jobs []*job.Job
	mu   sync.Mutex
}

func NewPriorityQueue() Queue {
	pq := &PriorityQueue{
		jobs: make([]*job.Job, 0),
	}
	heap.Init(pq)
	return pq
}

// heap.Interface implementation
func (pq *PriorityQueue) Len() int {
	return len(pq.jobs)
}
func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.jobs[i].Priority > pq.jobs[j].Priority
}
func (pq *PriorityQueue) Swap(i, j int) {
	pq.jobs[i], pq.jobs[j] = pq.jobs[j], pq.jobs[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	pq.jobs = append(pq.jobs, x.(*job.Job))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := pq.jobs
	n := len(old)
	item := old[n-1]
	pq.jobs = old[0 : n-1]
	return item
}

// Queue interface implementation - rename these instead
func (pq *PriorityQueue) Length() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.Len()
}

func (pq *PriorityQueue) isQueueJobsEmpty() bool {
	return pq.Length() == 0
}

func (pq *PriorityQueue) EnqueueJob(job *job.Job) error {
	log.Printf("Enqueueing job %s with priority %d", job.ID, job.Priority)
	if job == nil {
		return errors.New("cannot push nil job")
	}
	pq.mu.Lock()
	defer pq.mu.Unlock()
	heap.Push(pq, job)
	return nil
}

func (pq *PriorityQueue) DequeueJob() (*job.Job, error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.Len() == 0 {
		return nil, errors.New("queue is empty")
	}
	return heap.Pop(pq).(*job.Job), nil
}

func (pq *PriorityQueue) Peek() (*job.Job, error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.Len() == 0 {
		return nil, errors.New("queue is empty")
	}
	return pq.jobs[0], nil
}
