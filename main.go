package main

import (
	"github.com/lizongying/cron/cron"
	"log"
	"os"
	"time"
)

func main() {
	var logger cron.Logger
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	t := cron.New(cron.WithIntervalSecond(), cron.WithLoggerStdout())
	//t := cron.New()

	job := cron.Job{
		Spec:    "@every 3 seconds",
		RunType: cron.Divisibility,
		Id:      1,
		Meta:    nil,
		Callback: func(id int, meta any, now time.Time) {
			logger.Println(id, meta, now)
		},
	}
	_ = t.AddJob(&job)

	_ = t.Start()
	defer func() {
		_ = t.Stop()
	}()

	select {}
}
