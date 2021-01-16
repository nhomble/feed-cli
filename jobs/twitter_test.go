package jobs

import "testing"

const twitterUrl = "twitter://@foo"

func TestIsTwitter(t *testing.T) {
	job := createJob([]byte(twitterUrl))
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
	user := parseTwitterUser(createJob([]byte(twitterUrl)))
	if user != "foo" {
		t.Fail()
	}
}
