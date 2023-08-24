package cron

import (
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	t.Log("running", c.running)
}

func TestCron_AddJobCallbackNil(t *testing.T) {
	c := New()
	_, err := c.AddJob("every 2 seconds", nil)
	if err != nil {
		t.Log(err)
		return
	}
}

func TestCron_AddJobIdNil(t *testing.T) {
	c := New()
	_, err := c.AddJob("every 2 seconds", func() {
		t.Log(1)
	})
	if err != nil {
		t.Log(err)
		return
	}
}

func TestCron_AddJobExists(t *testing.T) {
	var err error
	c := New()
	_, err = c.AddJob("every 2 seconds", func() {
		t.Log(1)
	})
	if err != nil {
		t.Log(err)
		return
	}
	_, err = c.AddJob("every 2 seconds", func() {
		t.Log(1)
	})
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("%+v", c)
}

func TestCron_AddJob(t *testing.T) {
	c := New()
	_, err := c.AddJob("every 2 seconds", func() {
		t.Log(1)
	})
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("%+v", c)
}

func TestCron_StopWhenNotRunning(t *testing.T) {
	c := New(WithSecond(), WithStdout())
	c.MustStop()
}

func TestCron_Stop(t *testing.T) {
	c := New(WithSecond(), WithStdout())
	c.MustStart()
	c.MustStop()
	t.Log("running", c.running)
}

func TestCron_StartWhenRunning(t *testing.T) {
	c := New(WithSecond(), WithStdout())
	c.MustStart()
	c.MustStart()
}

func TestCron_ReStart(t *testing.T) {
	c := New(WithSecond(), WithStdout())
	c.MustStart()
	c.MustStop()
	c.MustStart()
}
