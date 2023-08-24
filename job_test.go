package cron

import (
	"testing"
	"time"
)

func TestJob_Init(t *testing.T) {
	now := time.Now()
	t.Log(now)
	job := new(Job).
		EverySecond(10).
		MustSince("10:15").
		Callback(func() {
			t.Log("callback")
		})
	if err := job.init(uint32(now.Unix())); err != nil {
		t.Log(err)
		return
	}

	t.Log(job.NextTime())
}

func TestJob_Next(t *testing.T) {
	now := time.Now()
	t.Log(now)
	job := new(Job).
		EverySecond(10).
		MustSince("10:15").
		Callback(func() {
			t.Log("callback")
		})
	if err := job.init(uint32(now.Unix())); err != nil {
		t.Log(err)
		return
	}
	t.Log(job.NextTime())

	job.next()

	t.Log(job.NextTime())
}
