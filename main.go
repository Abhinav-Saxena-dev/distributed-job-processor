package main

import (
	"distributed-job-processor/api"
	"distributed-job-processor/queue"
	"distributed-job-processor/worker"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// In main.go
func main() {
	var jobQueue queue.Queue
	usePriorityQueue := true

	if usePriorityQueue {
		jobQueue = queue.NewPriorityQueue()
	} else {
		jobQueue = queue.NewQueue()
	}

	// Create worker pool
	pool := worker.NewWorkerPool(5, jobQueue)
	pool.Start()

	// Create and start HTTP server
	server := api.NewJobServer(pool)

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Start("8080"); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	log.Println("Shutting down...")

	// Graceful shutdown
	if err := server.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}
