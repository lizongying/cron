package cron

import (
	"testing"
)

func TestWithSecond(t *testing.T) {
	c := New(WithSecond())
	t.Log("interval", c.interval)
}

func TestWithMinute(t *testing.T) {
	c := New(WithMinute())
	t.Log("interval", c.interval)
}

func TestWithLogger(t *testing.T) {
	logger := NewLoggerStdout()
	c := New(WithLogger(logger))
	t.Log("logger", c.logger)
}

func TestWithStdout(t *testing.T) {
	c := New(WithStdout())
	t.Log("logger", c.logger)
}

func TestWithDivisibility(t *testing.T) {
	c := New(WithDivisibility())
	t.Log("divisibility", c.divisibility)
}
