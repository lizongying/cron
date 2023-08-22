package main

import (
	"github.com/lizongying/cron/cron"
	"time"
)

func main() {
	logger := cron.NewLoggerStdout()
	c := cron.New(cron.WithSecond())
	for i := 0; i < 1000000; i++ {
		v := i
		c.MustAddJob("*/3 * * * * *", &cron.Job{
			Callback: func() {
				if (v+1)%1000000 == 0 {
					logger.Info(v, time.Now())
				}
			},
		})
	}
	logger.Info("now", time.Now())
	c.MustStart()
	select {}
}
