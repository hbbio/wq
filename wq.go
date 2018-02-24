// (c) Henri Binsztok, 2015
// See: LICENSE

package wq

import (
	"database/sql"
	"fmt"
	"gopkg.in/cheggaaa/pb.v1"
	"log"
	"os"
	"time"
)

var waitBeforeStart = true
var verbose = false

func SetWait(b bool) {
	waitBeforeStart = b
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func printJob(nb int) {
	plural := ""
	if nb > 1 {
		plural = "s"
	}
	fmt.Print("running ", nb, " job", plural) // , list)
	for i := 1; i <= 3; i++ {
		if waitBeforeStart {
			time.Sleep(1000 * time.Millisecond)
		}
		fmt.Print(".")
	}
	fmt.Println("")
}

// payload do not have results
type payload func(*sql.DB, string) error

func worker(db *sql.DB, id int, fn payload, jobs <-chan string, done chan<- bool) {
	var err error
	for job := range jobs {
		if verbose {
			fmt.Println("worker", id, ":", job)
		}
		err = fn(db, job)
		if err != nil {
			log.Println(os.Stderr, err)
		}
		done <- err == nil
	}
}

func Queue(db *sql.DB, fn payload, nbWorkers int, list []string) int {
	// compute and check number of jobs > 0
	nbJobs := len(list)
	if nbJobs == 0 || nbWorkers <= 0 {
		fmt.Println("nothing to do, quitting")
		return 0
	}
	printJob(nbJobs)

	bar := pb.StartNew(nbJobs)

	// compute pool size
	poolSize := Min(nbJobs, nbWorkers)

	// make channels
	jobs := make(chan string, poolSize)
	done := make(chan bool, poolSize)
	var result bool
	correct := 0

	// 1. create pool of workers
	for w := 1; w <= poolSize; w++ {
		go worker(db, w, fn, jobs, done)
	}

	// 2. iterate, feed initial pool
	i := 0
	for i < poolSize {
		jobs <- list[i]
		i++
	}

	// 3. iterate, feed pool when a job is finished
	for i < nbJobs {
		result = <-done // wait for one job to finish
		bar.Increment()
		if result {
			correct++
		}
		jobs <- list[i] // launch a new job
		i++
	}
	close(jobs)

	// 4. wait for remaining workers to finish
	i = 0
	for i < poolSize {
		result = <-done
		bar.Increment()
		if result {
			correct++
		}
		i++
	}
	close(done) // necessary?

	if correct < nbJobs {
		bar.FinishPrint(fmt.Sprintf("⚠ some jobs had errors: %d/%d", correct, nbJobs))
	} else {
		bar.FinishPrint("✓ ok")
	}
	return correct
}
