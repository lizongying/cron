package cron

import (
	"testing"
)

func TestWithSecond(t *testing.T) {
	tw := New(WithIntervalSecond())
	t.Log("slotCount", tw.slotCount)
}

func TestWithMinute(t *testing.T) {
	tw := New(WithIntervalMinute())
	t.Log("slotCount", tw.slotCount)
}

func TestWithLogger(t *testing.T) {
	logger := NewLoggerStdout()
	tw := New(WithLogger(logger))
	t.Log("logger", tw.logger)
}

func TestWithLoggerStdout(t *testing.T) {
	tw := New(WithLoggerStdout())
	t.Log("logger", tw.logger)
}
