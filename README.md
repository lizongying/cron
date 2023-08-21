# Cron

基于时间轮实现的定时任务，更好的性能。目前支持crontab格式或`every 1 second|minute|hour|day|month|week`

[cron](https://github.com/lizongying/cron)

## Features

* 支持crontab格式或`every 1 second|minute|hour|day|month|week`，更简单
* 执行时间进行了修正，会在秒/分开始的时候才执行，所以初次执行会有不到1秒/1分的延时
* 基于时间轮，保证实时性，任务容量会更高些。并发性能是github.com/robfig/cron的百倍以上，且更准时
* 回调函数可以增加一些参数，更容易调试，使用也更方便
* 支持整时执行和立即执行

## Install

```shell
go get -u github.com/lizongying/cron
```

## Usage

* Spec: 定时
* OnlyOnce: 只执行一次
* RunIfDelay: 即使超时(超过最大job处理数量)也会执行，否则本次不执行。
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
		//Spec:     "every 3 seconds",
		Spec:     "*/3 * * * * *",
		OnlyOnce: false,
		RunType:  cron.Divisibility,
		Id:       1,
		Callback: func(id int, meta any, t time.Time) {
			logger.Println(id, meta, t)
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

## Tips

* 建议秒级别最大任务控制在1,000,000(Apple M1 Pro, 32 GB))以内，防止任务超时。可能支持更大数量，请自行测试。

## TODO
