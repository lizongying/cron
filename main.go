package main

import (
	"github.com/lizongying/cron/cron"
	"log"
	"os"
)

func main() {
	var logger cron.Logger
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	t := cron.New(cron.WithIntervalSecond(), cron.WithLoggerStdout())

	job := cron.Job{
		Spec:     "@every 3 seconds",
		OnlyOnce: false,
		RunType:  cron.Divisibility,
		Id:       1,
		Callback: func() {
			logger.Println(1)
		},
	}
	_ = t.AddJob(&job)

	_ = t.Start()
	defer func() {
		_ = t.Stop()
	}()

	select {}
}
