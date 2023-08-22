package cron

import (
	"testing"
)

func TestWithSecond(t *testing.T) {
	tw := New(WithSecond())
	t.Log("slotCount", tw.slotCount)
}

func TestWithMinute(t *testing.T) {
	tw := New(WithMinute())
	t.Log("slotCount", tw.slotCount)
}

func TestWithLogger(t *testing.T) {
	logger := NewLoggerStdout()
	tw := New(WithLogger(logger))
	t.Log("logger", tw.logger)
}

func TestWithStdout(t *testing.T) {
	tw := New(WithStdout())
	t.Log("logger", tw.logger)
}
