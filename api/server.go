package api

import (
	"distributed-job-processor/job"
	"distributed-job-processor/worker"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type JobServer struct {
	workerPool *worker.WorkerPool
}

func NewJobServer(wp *worker.WorkerPool) *JobServer {
	return &JobServer{
		workerPool: wp,
	}
}

func (s *JobServer) Start(port string) error {
	// Register routes
	http.HandleFunc("/jobs", s.handleJobs)

	log.Printf("Starting server on port %s", port)
	return http.ListenAndServe(":"+port, nil)
}

func (s *JobServer) handleJobs(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Type    job.JobType `json:"type"`
		Payload interface{} `json:"payload"`
	}

	// Decode JSON from request body into our struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Request recieved with type %s", req.Type)

	// Create new job
	newJob := job.NewJob(generateID(), req.Type, req.Payload)

	// Submit job using WorkerPool method
	if err := s.workerPool.SubmitJob(newJob); err != nil {
		http.Error(w, "Failed to queue job", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"job_id": newJob.ID,
	})
}

func (s *JobServer) Stop() error {
	s.workerPool.Stop()
	return nil
}

// Helper function to generate job ID
func generateID() string {
	// For now, just use timestamp
	return fmt.Sprintf("job-%d", time.Now().UnixNano())
}
