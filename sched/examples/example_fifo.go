package main

import (
	"log"

	"github.com/althk/goeasy/sched"
)

func main() {

	// All scheduler implementations adhere to TaskScheduler interface
	var fifoSched sched.TaskScheduler

	// The maximum number to tasks to be allowed before blocking
	var maxTasks uint = 10

	// The maximum number of tasks that will be run concurrently
	// Setting to 3 means at most 3 tasks will be running at any
	// given time
	var maxConcurrency uint = 3

	// Create a new FIFO Scheduler
	// FIFO Scheduler is threadsafe and supports all operations
	// in a concurrent environment
	fifoSched = sched.NewFIFO(maxTasks, maxConcurrency)

	// Initialize the scheduler
	if err := fifoSched.Init(); err != nil {
		log.Fatalf("scheduler init failed: %v", err)
	}

	// Schedule some tasks
	for i := 0; i < 5; i++ {
		// EnqueueTask takes a nullary function (function with no args)
		if err := fifoSched.EnqueueTask(func() {
			// do some work
		}); err != nil {
			log.Printf("Failed to enqueue task %v", err)
		}
	}

	// Shutdown the scheduler
	// Passing `true` to Shutdown() will ensure the scheduler
	// runs all the tasks already enqueued and then exit.
	// If `false` is passed, then the scheduler will finish
	// in-flight tasks and exit.
	if err := fifoSched.Shutdown(true); err != nil {
		log.Printf("Failed to shutdown %v", err)
	}

}
