package cron

import (
	"log"
	"os"
)

type Logger interface {
	Info(...any)
	Infof(format string, v ...any)
	Error(...any)
	Errorf(format string, v ...any)
}

type LoggerNothing struct{}

func (l *LoggerNothing) Info(_ ...any) {}

func (l *LoggerNothing) Infof(_ string, _ ...any) {}

func (l *LoggerNothing) Error(_ ...any) {}

func (l *LoggerNothing) Errorf(_ string, _ ...any) {}

type LoggerStdout struct {
	logger *log.Logger
}

func NewLoggerStdout() Logger {
	return &LoggerStdout{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile),
	}
}

func (l *LoggerStdout) Info(v ...any) {
	l.logger.Println(v...)
}

func (l *LoggerStdout) Infof(format string, v ...any) {
	l.logger.Printf(format, v...)
}

func (l *LoggerStdout) Error(v ...any) {
	l.logger.Println(v...)
}

func (l *LoggerStdout) Errorf(format string, v ...any) {
	l.logger.Printf(format, v...)
}
