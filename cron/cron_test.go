package cron

import (
	"testing"
)

func TestNew(t *testing.T) {
	tw := New()
	t.Log("slotCount", tw.slotCount)
}

func TestCron_AddJobCallbackNil(t *testing.T) {
	tw := New()
	job := Job{}
	_, err := tw.AddJob("every 2 seconds", &job)
	if err != nil {
		t.Log(err)
		return
	}
}

func TestCron_AddJobIdNil(t *testing.T) {
	tw := New()
	job := Job{
		Callback: func() {
			t.Log(1)
		},
	}
	_, err := tw.AddJob("every 2 seconds", &job)
	if err != nil {
		t.Log(err)
		return
	}
}

func TestCron_AddJobExists(t *testing.T) {
	var err error
	tw := New()
	job := Job{
		Callback: func() {
			t.Log(1)
		},
	}
	_, err = tw.AddJob("every 2 seconds", &job)
	if err != nil {
		t.Log(err)
		return
	}
	_, err = tw.AddJob("every 2 seconds", &job)
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("%+v", tw)
}

func TestCron_AddJob(t *testing.T) {
	tw := New()
	job := Job{
		Callback: func() {
			t.Log(1)
		},
	}
	_, err := tw.AddJob("every 2 seconds", &job)
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
	tw := New(WithMinute())
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
	t.Log("slotCount", tw.slotCount)
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
