package job

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type JobType string
type JobStatus string
type Priority int

const (
	LowPriority    Priority = 1
	MediumPriority Priority = 2
	HighPriority   Priority = 3
)

const (
	EmailJob    JobType = "email"
	APIJob      JobType = "api"
	DataJobType JobType = "data"
)

// Add status constants
const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type JobProgress struct {
	Status     JobStatus
	Progress   float64
	StartedAt  time.Time
	FinishedAt time.Time
}

type Job struct {
	ID          string
	Type        JobType
	Priority    Priority
	Payload     interface{}
	JobProgress JobProgress
	CreatedAt   time.Time
}

func (j *Job) String() string { // *Job means we will be directly effecting the original struct.
	return fmt.Sprintf("Job{ID: %s, Type: %s, Status: %s, CreatedAt: %s}",
		j.ID, j.Type, j.JobProgress.Status, j.CreatedAt)
}

func (j *Job) Validate() error {
	if j.ID == "" {
		return errors.New("job ID cannot be empty")
	}
	if j.Type == "" {
		return errors.New("job type cannot be empty")
	}
	if j.Payload == nil {
		return errors.New("job payload cannot be nil")
	}
	return nil
}

func NewJob(id string, jobType JobType, priority Priority, payload interface{}) *Job {
	jobProgress := JobProgress{
		Status: StatusPending,
	}
	log.Printf("Created new Job of type %s - ID %s", jobType, id)
	return &Job{
		ID:          id,
		Type:        jobType,
		Payload:     payload,
		Priority:    priority,
		JobProgress: jobProgress,
		CreatedAt:   time.Now(),
	}
}
