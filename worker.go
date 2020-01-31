package wq

import (
	"fmt"
	"log"
)

// Worker is a task executor.
type Worker struct {
	ID  int  // worker ID
	job *Job // pointer to Job
}

// NewWorker creates a new Worker.
func (j *Job) NewWorker(ID int) *Worker {
	w := Worker{ID: ID, job: j}
	return &w
}

// Launch launches a Worker.
func (w *Worker) Launch(tasks <-chan string, done chan<- bool) {
	for task := range tasks {
		if w.job.conf.Verbose {
			fmt.Printf("worker [%d]: task %s\n", w.ID, task)
		}
		err := w.job.fn(task)
		if err != nil {
			log.Println(err)
		}
		done <- err == nil
	}
}
