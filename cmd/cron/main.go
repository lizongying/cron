package main

import (
	"github.com/lizongying/cron"
)

func main() {
	logger := cron.NewLoggerStdout()
	c := cron.New(cron.WithFix())
	id := c.MustAddJob(new(cron.Job).
		EverySecond(3).
		Callback(func() {
			logger.Info("callback")
		}))
	logger.Info("id", id)
	c.MustStart()
	select {}
}
