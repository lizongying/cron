package main

import (
	"github.com/lizongying/cron/cron"
	"log"
	"runtime"
	"time"
)

func main() {
	begin := time.Now()
	c := cron.New(cron.WithSecond())
	for i := 0; i < 5000000; i++ {
		v := i
		_ = c.MustAddJob("every 3 second", func() {
			if v%5000000 == 0 {
				now := time.Now()
				log.Println(v, now.Sub(begin))
				begin = now
				var mem runtime.MemStats
				runtime.ReadMemStats(&mem)
				log.Printf("TotalAlloc = %v MiB\n", mem.TotalAlloc/1024/1024)
			}
		})
	}
	c.MustStart()
	log.Println("now", begin)
	select {}
}
