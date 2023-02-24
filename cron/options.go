package cron

import (
	"log"
	"os"
	"time"
)

type Options func(t *cron)

func WithIntervalSecond() Options {
	return func(t *cron) {
		t.interval = time.Second
		t.slotCount = 1 << 29
	}
}

func WithIntervalMinute() Options {
	return func(t *cron) {
		t.interval = time.Minute
		t.slotCount = 1 << 23
	}
}

func WithLogger(logger Logger) Options {
	return func(t *cron) {
		t.logger = logger
	}
}

func WithLoggerStdout() Options {
	return func(t *cron) {
		var logger Logger
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
		t.logger = logger
	}
}
