package main

import (
	"fmt"
	"github.com/lizongying/cron"
	"log"
	"runtime"
	"time"
)

func main() {
	logger := cron.NewLoggerStdout()
	c := cron.New(cron.WithStdout())
	c.MustStart()
	id := c.MustAddJob(new(cron.Job).
		MustEverySpec("3s").
		Callback(func() {
			fmt.Println(time.Now())
		}))
	logger.Info("id", id)
	select {}
}

func main1() {
	num := 4000000
	begin := time.Now()
	c := cron.New()
	c.MustStart()
	for i := 1; i <= num; i++ {
		v := i
		_ = c.MustAddJob(new(cron.Job).
			EverySecond(3).
			Callback(func() {
				if v == num {
					var mem runtime.MemStats
					runtime.ReadMemStats(&mem)
					now := time.Now()
					log.Printf("Alloc = %v M, Spend = %v\n", mem.Alloc/1024/1024, now.Sub(begin))
					begin = now
				}
			}))
	}
	log.Println("num", num)
	select {}
}
