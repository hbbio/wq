package wq

// Config describes the configuration of a job queue.
type Config struct {
	Workers         int  // number of workers
	Verbose         bool // verbose flag
	WaitBeforeStart bool // wait 3 seconds before starting
	ProgressBar     bool // display progress bar
}

// NewConfig creates a new configuration by default.
func NewConfig(workers int) *Config {
	return &Config{Workers: workers, ProgressBar: true}
}
