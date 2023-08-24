package cron

import (
	"time"
)

type Options func(t *Cron)

func WithSecond() Options {
	return func(t *Cron) {
		t.interval = time.Second
	}
}

func WithMinute() Options {
	return func(t *Cron) {
		t.interval = time.Minute
	}
}

func WithLogger(logger Logger) Options {
	return func(t *Cron) {
		t.logger = logger
	}
}

func WithStdout() Options {
	return func(t *Cron) {
		t.logger = NewLoggerStdout()
	}
}

func WithDivisibility() Options {
	return func(t *Cron) {
		t.divisibility = true
	}
}
