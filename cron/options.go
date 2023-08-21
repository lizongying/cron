package cron

import (
	"time"
)

type Options func(t *Cron)

func WithIntervalSecond() Options {
	return func(t *Cron) {
		t.interval = time.Second
		t.slotCount = 60 * 60 * 24 * 366
	}
}

func WithIntervalMinute() Options {
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

func WithLoggerStdout() Options {
	return func(t *Cron) {
		t.logger = NewLoggerStdout()
	}
}
