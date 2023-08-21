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

	t := cron.New(cron.WithIntervalSecond())

	for i := 0; i < 5000000; i++ {
		v := i
		t.MustAddJob(&cron.Job{
			Spec: "*/3 * * * * *",
			Meta: v,
			Id:   v + 1,
			Callback: func(id int, meta any, time time.Time) {
				if id%5000000 == 0 {
					logger.Println(id, meta, time)
				}
			},
		})
	}

	logger.Println(time.Now())
	t.MustStart()
	defer func() {
		t.MustStop()
	}()

	select {}
}
