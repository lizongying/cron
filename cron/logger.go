package cron

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
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
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (l *LoggerStdout) Info(v ...any) {
	_, file, line, _ := runtime.Caller(1)
	v = append([]any{strings.Join([]string{file, strconv.Itoa(line)}, ":")}, v...)
	l.logger.Println(v...)
}

func (l *LoggerStdout) Infof(format string, v ...any) {
	_, file, line, _ := runtime.Caller(1)
	format = fmt.Sprintf("%s %s", strings.Join([]string{file, strconv.Itoa(line)}, ":"), format)
	l.logger.Printf(format, v...)
}

func (l *LoggerStdout) Error(v ...any) {
	_, file, line, _ := runtime.Caller(1)
	v = append([]any{strings.Join([]string{file, strconv.Itoa(line)}, ":")}, v...)
	l.logger.Println(v...)
}

func (l *LoggerStdout) Errorf(format string, v ...any) {
	_, file, line, _ := runtime.Caller(1)
	format = fmt.Sprintf("%s %s", strings.Join([]string{file, strconv.Itoa(line)}, ":"), format)
	l.logger.Printf(format, v...)
}
