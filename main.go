package main

import (
	"github.com/lizongying/cron/cron"
	"log"
	"os"
)

func main() {
	var logger cron.Logger
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	//t := cron.NewTimeWheel(cron.WithIntervalSecond(), cron.WithLoggerStdout())
	t := cron.NewCron()

	job := cron.Job{
		Spec: "@every 1 minutes",
		Id:   1,
		Meta: nil,
		Callback: func(id int, meta any) {
			logger.Println(id, meta)
		},
	}
	_ = t.AddJob(&job)

	_ = t.Start()
	defer func() {
		_ = t.Stop()
	}()

	select {}
}
