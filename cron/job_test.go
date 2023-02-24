package cron

import (
	"testing"
	"time"
)

func TestJob_Next(t *testing.T) {
	t.Log(time.Now())
	job := Job{
		Spec: "@every 2 seconds",
	}
	err := job.Next(time.Second)
	if err != nil {
		t.Log(err)
	}
	t.Logf("%+v", job)
	t.Logf("%+v", job.slot)
	t.Log(job.nextTime)
}

func TestJob_GetSlot(t *testing.T) {
	now, _ := time.ParseInLocation(time.DateTime, "2023-12-31 23:59:59", time.Local)
	slot := GetSlot(now, time.Minute)
	t.Log(slot)
	slot = GetSlot(now, time.Second)
	t.Log(slot)
}
