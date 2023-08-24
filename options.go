package cron

type Options func(t *Cron)

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
