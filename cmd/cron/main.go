package main

import (
	"github.com/lizongying/cron"
	"log"
)

func main() {
	c := cron.New(cron.WithSecond())
	id := c.MustAddJob("every 3 second", func() {
		log.Println("callback")
	})
	log.Println("id", id)
	c.MustStart()
	select {}
}
