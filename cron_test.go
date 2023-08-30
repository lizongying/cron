package cron

import (
	"github.com/lizongying/gooptimizer"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c := New(WithStdout())
	t.Log("running", c.running)
}

func TestCron_AddJob(t *testing.T) {
	t.Log(time.Now())
	c := New(WithStdout())
	id, err := c.AddJob(new(Job).
		EverySecond(10).
		MustSince("10:15").
		Callback(func() {
			t.Log("callback")
		}))

	if err != nil {
		t.Log(err)
		return
	}

	t.Log(id)
	t.Log(c.MustGetJob(id).NextTime())
}

func TestCron_AddJobNil(t *testing.T) {
	c := New(WithStdout())
	_, err := c.AddJob(nil)
	if err != nil {
		t.Log(err)
		return
	}
}

func TestCron_StopWhenNotRunning(t *testing.T) {
	c := New(WithStdout())
	c.MustStop()
}

func TestCron_Stop(t *testing.T) {
	c := New(WithStdout())
	c.MustStart()
	c.MustStop()
}

func TestCron_StartWhenRunning(t *testing.T) {
	c := New(WithStdout())
	c.MustStart()
	c.MustStart()
}

func TestCron_Restart(t *testing.T) {
	c := New(WithStdout())
	c.MustStart()
	c.MustStop()
	c.MustStart()
}

func TestCron_Optimize(t *testing.T) {
	gooptimizer.StructAlignWithPrint(new(Cron))
}
