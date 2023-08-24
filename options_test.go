package cron

import (
	"testing"
)

func TestWithLogger(t *testing.T) {
	logger := NewLoggerStdout()
	c := New(WithLogger(logger))
	t.Log("logger", c.logger)
}

func TestWithStdout(t *testing.T) {
	c := New(WithStdout())
	t.Log("logger", c.logger)
}
