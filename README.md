# Cronjob

基于时间轮实现的定时任务，目前支持 @every 1 second|minute|hour|day|month|week

# Install

```shell
go get -u github.com/lizongying/cron
```

# Usage

* Spec: 定时
* OnlyOnce: 只执行一次
* RunType: now 基于当前时间立即执行; Divisibility 整时运行

```go
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

	job := cron.Job{
		Spec:     "@every 3 seconds",
		OnlyOnce: false,
		RunType:  cron.Divisibility,
		Id:       1,
		Meta:     nil,
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
```

# Compare

|     | lizongying/cronjob                         | robfig/cron                       |
|-----|--------------------------------------------|-----------------------------------|
|     | 暂未支持crontab格式。使用自定义格式，更简单                  | 支持crontab格式，使用更广泛                 |
|     | 管理任务和运行任务不同线程，任务运行更实时                      | 管理任务和运行任务相同线程，在任务密集操作时，可能会有影响任务运行 |
|     | 执行时间进行了修正，会在秒/分开始的时候才执行，所以初次执行会有不到1秒/1分的延时 | 基于当前时间执行                          |
|     | 基于时间轮，保证实时性，任务容量会更高些                       | 基于队列，可能会有延时，任务多时会影响比较大            |
|     | 回调函数增加额外参数（id、meta、time）， 使用更方便            | 无。更简洁。                            |
|     | 无回调函数的处理方法，交给调度者处理                         | 回调函数有些处理方法，可以直接使用。                |
|     | 支持整点执行和立即执行                                | 可能麻烦些                             |