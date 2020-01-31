package wq

import (
	"errors"
	"fmt"
	"time"

	"github.com/schollz/progressbar"
)

// Payload defines the payload function.
//
// Payloads do not return value and rely on side effects (updating a database, printing values, etc.).
type Payload func(string) error

// Job defines a list of tasks to execute.
type Job struct {
	conf  *Config
	fn    Payload
	tasks []string
}

// ErrNoJob is the error used when there is nothing to do.
var ErrNoJob = errors.New("nothing to do")

// NewJob creates a new job.
func (conf *Config) NewJob(fn Payload, tasks []string) (*Job, error) {
	if len(tasks) == 0 || conf.Workers <= 0 {
		return nil, ErrNoJob
	}
	return &Job{conf: conf, fn: fn, tasks: tasks}, nil
}

// Print displays a job with potential confirmation.
func (j *Job) Print() {
	nb := len(j.tasks)
	fmt.Print("running ", nb, " task", Plural(nb)) // , list)
	if j.conf.WaitBeforeStart {
		for i := 1; i <= 3; i++ {
			time.Sleep(1000 * time.Millisecond)
			fmt.Print(".")
		}
		fmt.Println("")
	}
}

// Run executes a job.
func (j *Job) Run() int {
	var bar *progressbar.ProgressBar

	nbTasks := len(j.tasks)

	if j.conf.ProgressBar {
		bar = progressbar.New(nbTasks)
	}

	// compute pool size
	poolSize := MinInt(nbTasks, j.conf.Workers)

	// make channels
	chanTasks := make(chan string, poolSize)
	defer close(chanTasks)
	chanDone := make(chan bool, poolSize)
	defer close(chanDone)

	// 1. create pool of workers
	for w := 1; w <= poolSize; w++ {
		worker := j.NewWorker(w)
		go worker.Launch(chanTasks, chanDone)
		// go conf.worker(w, fn, chanTasks, done)
	}

	// 2. iterate, feed initial pool
	for i := 0; i < poolSize; i++ {
		chanTasks <- j.tasks[i]
	}

	correct := 0
	// 3. iterate, feed pool when a task is finished
	for i := poolSize; i < nbTasks; i++ {
		result := <-chanDone // wait for one task to finish
		if result {
			correct++
		}
		chanTasks <- j.tasks[i] // launch a new task
		if j.conf.ProgressBar {
			bar.Add(1)
		}
	}

	// 4. wait for remaining workers to finish
	for i := 0; i < poolSize; i++ {
		result := <-chanDone
		if result {
			correct++
		}
		if j.conf.ProgressBar {
			bar.Add(1)
		}
	}

	if j.conf.Verbose {
		if correct < nbTasks {
			fmt.Printf("⚠ some tasks had errors: %d/%d\n", correct, nbTasks)
		} else {
			fmt.Println("✓ ok")
		}
	}
	return correct
}
