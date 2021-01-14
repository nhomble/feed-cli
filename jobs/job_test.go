package jobs

import (
	"testing"
	"time"
)

const url = "https://my.feed.com/feed"

func TestDefaultJob(t *testing.T) {
	job := createJob([]byte(url))
	if job.link != url {
		t.Fail()
	}
}

func TestTrimUrl(t *testing.T) {
	job := createJob([]byte("    " + url + "     "))
	if job.link != url {
		t.Fail()
	}
}

func TestParseLimit(t *testing.T) {
	job := createJob([]byte(url + "   limit=20"))
	if job.limit != 20 {
		t.Fail()
	}
}

func TestParseDaysOld(t *testing.T) {
	job := createJob([]byte(url + "   daysOld=9"))
	if job.age != 9*24*time.Hour {
		t.Fail()
	}
}

func TestCanProcess_Something(t *testing.T) {
	if !canProcess([]byte("https://my.feed/feed")) {
		t.Fail()
	}
}

func TestCanProcess_Comment(t *testing.T) {
	if canProcess([]byte("  # some comment")) {
		t.Fail()
	}
}

func TestCanProcess_Empty(t *testing.T) {
	if canProcess([]byte{}) {
		t.Fail()
	}
}

func TestCanProcess_Whitespace(t *testing.T) {
	if canProcess([]byte("    ")) {
		t.Fail()
	}
}

func TestParseTimeout(t *testing.T) {
	job := createJob([]byte(url + "   timeout=100"))
	if job.timeout != time.Duration(100)*time.Second {
		t.Fail()
	}
}
