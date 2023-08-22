package cron

import (
	"testing"
	"time"
)

func TestJob_Init(t *testing.T) {
	now := time.Now()
	t.Log(now)
	job := &Job{}
	if err := job.Init("every 2 minutes", time.Minute); err != nil {
		t.Log(err)
		return
	}

	t.Log(job.nextTime)
}

func TestJob_InitCrontab(t *testing.T) {
	var err error
	now := time.Now()
	t.Log(now)
	job := &Job{}
	//err := job.Init("* * * * * *", time.Minute)
	//if err != nil {
	//	t.Log(err)
	//}
	//t.Logf("%+v", job)
	//t.Log(job.nextTime)

	job = &Job{}
	if err = job.Init("* 1-10 1,2,3 4 0-6", time.Minute); err != nil {
		t.Log(err)
		return
	}

	t.Log(job.nextTime)
	if err = job.Next(time.Minute); err != nil {
		t.Log(err)
	}
	t.Log(job.nextTime)
}

func TestJob_Next(t *testing.T) {
	var err error
	now := time.Now()
	t.Log(now)
	job := &Job{}
	if err = job.Init("every 2 seconds", time.Second); err != nil {
		t.Log(err)
		return
	}

	if err = job.Next(time.Second); err != nil {
		t.Log(err)
		return
	}

	t.Log("slot", job.Slot())
	t.Log(job.nextTime)
}

func TestJob_GetSlot(t *testing.T) {
	now, _ := time.ParseInLocation(time.DateTime, "2023-12-31 23:59:59", time.Local)
	slot := SlotSinceYear(now, time.Minute)
	t.Log(slot)
	slot = SlotSinceYear(now, time.Second)
	t.Log(slot)
}

func TestJob_SlotSinceYear(t *testing.T) {
	year, _ := time.ParseInLocation("2006", time.Now().Format("2006"), time.Local)
	now, _ := time.ParseInLocation(time.DateTime, "2023-01-01 00:00:00", time.Local)
	t.Log(year, now)
	slot := SlotSinceYear(now, time.Minute)
	t.Log(slot)
	slot = SlotSinceYear(now, time.Second)
	t.Log(slot)

	now, _ = time.ParseInLocation(time.DateTime, "2023-01-01 00:00:01", time.Local)
	t.Log(now)
	slot = SlotSinceYear(now, time.Minute)
	t.Log(slot)
	slot = SlotSinceYear(now, time.Second)
	t.Log(slot)

	now, _ = time.ParseInLocation(time.DateTime, "2023-01-01 00:00:59", time.Local)
	t.Log(now)
	slot = SlotSinceYear(now, time.Minute)
	t.Log(slot)
	slot = SlotSinceYear(now, time.Second)
	t.Log(slot)

	now, _ = time.ParseInLocation(time.DateTime, "2023-01-01 00:01:00", time.Local)
	t.Log(now)
	slot = SlotSinceYear(now, time.Minute)
	t.Log(slot)
	slot = SlotSinceYear(now, time.Second)
	t.Log(slot)

	now, _ = time.ParseInLocation(time.DateTime, "2023-01-01 00:59:00", time.Local)
	t.Log(now)
	slot = SlotSinceYear(now, time.Minute)
	t.Log(slot)
	slot = SlotSinceYear(now, time.Second)
	t.Log(slot)

	now, _ = time.ParseInLocation(time.DateTime, "2023-12-31 23:59:59", time.Local)
	t.Log(now)
	slot = SlotSinceYear(now, time.Minute)
	t.Log(slot)
	slot = SlotSinceYear(now, time.Second)
	t.Log(slot)

	now, _ = time.ParseInLocation(time.DateTime, "2024-01-01 00:00:00", time.Local)
	t.Log(now)
	slot = SlotSinceYear(now, time.Minute)
	t.Log(slot)
	slot = SlotSinceYear(now, time.Second)
	t.Log(slot)
}
