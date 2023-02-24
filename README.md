# Cronjob

基于时间轮实现的定时任务，目前支持 @every 1 second|minute|hour|day|month|week

# Install

```shell
go get -u github.com/lizongying/cron
```

# Usage

```go
package main

import (
	"github.com/lizongying/cron/cron"
	"log"
	"os"
)

func main() {
	var logger cron.Logger
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	t := cron.NewCron(cron.WithIntervalSecond(), cron.WithLoggerStdout())

	job := cron.Job{
		Spec: "@every 2 seconds",
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
```

