package cron

type Logger interface {
	Println(...any)
	Printf(format string, v ...any)
}

type LoggerNothing struct{}

func (l *LoggerNothing) Println(_ ...any) {}

func (l *LoggerNothing) Printf(_ string, _ ...any) {}
