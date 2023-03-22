package cron

import (
	"testing"
	"time"
)

func TestClock_New(t *testing.T) {
	var err error
	var a uint64
	var b uint64
	var c uint64
	var d uint64
	var e uint64
	var f uint64
	for _, v := range []uint8{0} {
		a |= 1 << v
	}
	for _, v := range []uint8{1, 2, 3, 4, 5, 6} {
		b |= 1 << v
	}
	for _, v := range []uint8{1, 2, 3, 4, 5, 6} {
		c |= 1 << v
	}
	for _, v := range []uint8{1, 2, 3, 4, 5, 6} {
		d |= 1 << v
	}
	for _, v := range []uint8{1, 2, 3, 4, 5, 6} {
		e |= 1 << v
	}
	for _, v := range []uint8{1, 2, 3, 4, 5, 6} {
		f |= 1 << v
	}
	clock, err := NewClock(time.Minute, a, b, c, d, e, f)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(clock.String()) // 22:58:00
	now, err := clock.NextWithWeek()
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(now) // 22:58:01
}
