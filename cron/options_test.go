package cron

import (
	"log"
	"os"
	"testing"
)

func TestWithSecond(t *testing.T) {
	tw := Newcron(WithIntervalSecond())
	t.Logf("%+v", tw)
}

func TestWithMinute(t *testing.T) {
	tw := Newcron(WithIntervalMinute())
	t.Logf("%+v", tw)
}

func TestWithLogger(t *testing.T) {
	var logger Logger
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	tw := Newcron(WithLogger(logger))
	t.Logf("%+v", tw)
}
