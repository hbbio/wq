package wq

// Queue creates and runs a job queue.
//
// This all-in-one function does not require defining elements separately.
func Queue(workers int, fn Payload, tasks []string) (int, error) {
	conf := NewConfig(workers)
	job, err := conf.NewJob(fn, tasks)
	if err != nil {
		return 0, err
	}
	job.Print()
	return job.Run(), nil
}
