package cron

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tw := New()
	t.Logf("%+v", tw)
}

func TestCron_AddJobCallbackNil(t *testing.T) {
	tw := New()
	job := Job{
		Spec: "@every 2 seconds",
		Id:   0,
		Meta: nil,
	}
	err := tw.AddJob(&job)
	if err != nil {
		t.Log(err)
		return
	}
}

func TestCron_AddJobIdNil(t *testing.T) {
	tw := New()
	job := Job{
		Spec: "@every 2 seconds",
		Meta: nil,
		Callback: func(id int, meta any, now time.Time) {
			t.Log(id, meta, now)
		},
	}
	err := tw.AddJob(&job)
	if err != nil {
		t.Log(err)
		return
	}
}

func TestCron_AddJobExists(t *testing.T) {
	var err error
	tw := New()
	job := Job{
		Spec: "@every 2 seconds",
		Id:   0,
		Meta: nil,
		Callback: func(id int, meta any, now time.Time) {
			t.Log(id, meta, now)
		},
	}
	err = tw.AddJob(&job)
	if err != nil {
		t.Log(err)
		return
	}
	err = tw.AddJob(&job)
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("%+v", tw)
}

func TestCron_AddJob(t *testing.T) {
	tw := New()
	job := Job{
		Spec: "@every 2 seconds",
		Id:   0,
		Meta: nil,
		Callback: func(id int, meta any, now time.Time) {
			t.Log(id, meta, now)
		},
	}
	err := tw.AddJob(&job)
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("%+v", tw)
}

func TestCron_StopWhenNotRunning(t *testing.T) {
	tw := New()
	err := tw.Stop()
	if err != nil {
		t.Log(err)
		return
	}
}

func TestCron_Stop(t *testing.T) {
	var err error
	tw := New()
	err = tw.Start()
	if err != nil {
		t.Log(err)
		return
	}
	err = tw.Stop()
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("%+v", tw)
}

func TestCron_StartWhenRunning(t *testing.T) {
	var err error
	tw := New()
	err = tw.Start()
	if err != nil {
		t.Log(err)
		return
	}
	err = tw.Start()
	if err != nil {
		t.Log(err)
		return
	}
}

func TestCron_Start(t *testing.T) {
	var err error
	tw := New()
	err = tw.Start()
	if err != nil {
		t.Log(err)
		return
	}
	err = tw.Stop()
	if err != nil {
		t.Log(err)
		return
	}
	err = tw.Start()
	if err != nil {
		t.Log(err)
		return
	}
}
