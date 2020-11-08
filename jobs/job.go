package jobs

import (
	"strings"
	"time"
)

const DAYS = 24 * time.Hour

type Job struct {
	link    string
	limit   int
	age     time.Duration
	timeout time.Duration
}

func createJob(bytes []byte) Job {
	line := string(bytes)
	s := strings.TrimSpace(line)
	return Job{
		link:    s,
		limit:   5,
		age:     14 * DAYS,
		timeout: time.Minute,
	}
}
