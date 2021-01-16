package jobs

import "testing"

const twitterURL = "twitter://@foo"

func TestIsTwitter(t *testing.T) {
	job := createJob([]byte(twitterURL))
	if !job.isTwitter() {
		t.Fail()
	}
}

func TestIsNotTwitter(t *testing.T) {
	job := createJob([]byte(url))
	if job.isTwitter() {
		t.Fail()
	}
}

func TestParseHandle(t *testing.T) {
	user := parseTwitterUser(createJob([]byte(twitterURL)))
	if user != "foo" {
		t.Fail()
	}
}
