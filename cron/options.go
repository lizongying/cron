package cron

import (
	"time"
)

type Options func(t *Cron)

func WithSecond() Options {
	return func(t *Cron) {
		t.interval = time.Second
		t.slotCount = 60 * 60 * 24 * 366
	}
}

func WithMinute() Options {
	return func(t *Cron) {
		t.interval = time.Minute
		t.slotCount = 60 * 24 * 366
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
