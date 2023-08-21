package main

import (
	"github.com/lizongying/cron/cron"
	"time"
)

func main() {
	logger := cron.NewLoggerStdout()
	c := cron.New(cron.WithIntervalSecond())
	for i := 0; i < 500000; i++ {
		v := i
		c.MustAddJob(&cron.Job{
			Spec: "*/3 * * * * *",
			Meta: v,
			Id:   v + 1,
			Callback: func(id int, meta any) {
				if id%500000 == 0 {
					logger.Info(id, meta, time.Now())
				}
			},
		})
	}
	logger.Info("now", time.Now())
	c.MustStart()

	select {}
}
