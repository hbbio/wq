# wq, A basic Worker Queue in Go

# Example

```go
fn := func(key string) error { ... }
nb, err := wq.Queue(
    8, // number of workers
    fn, // payload function
    []string{"key1", ...},
)
```

# Detailed usage

```go
conf := wq.NewConfig(workers)
job, err := conf.NewJob(fn, tasks)
if err != nil {
	...
}
job.Print()
ok := job.Run() // ok is the number of successful tasks
...
```

# About

This project was originally written in 2015 while playing with Go.
It was a very simple worker/queue where the payload for each job will update state in a SQL database, later rewritten without the `database/sql` dependency.
