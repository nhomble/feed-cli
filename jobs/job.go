package jobs

import (
	"strconv"
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
	tokens := strings.Split(s, " ")

	job := Job{
		link:    tokens[0],
		limit:   5,
		age:     14 * DAYS,
		timeout: time.Minute,
	}

	for i := 1; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i], "limit=") {
			v := strings.TrimPrefix(tokens[i], "limit=")
			limit, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				panic(err)
			} else if i < 1 {
				panic("Must provide a positive limit")
			}
			job.limit = int(limit)
		} else if strings.HasPrefix(tokens[i], "daysOld=") {
			v := strings.TrimPrefix(tokens[i], "daysOld=")
			age, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				panic(err)
			} else if age < 1 {
				panic("Must provde a positive limit")
			}
			job.age = time.Duration(age*24) * time.Hour
		}
	}

	return job
}
