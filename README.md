# Cron

基于时间轮实现的定时任务，更好的性能。目前支持crontab格式或如`every 1 second|minute|hour|day|month|week`格式

[cron](https://github.com/lizongying/cron)

## Features

* 支持crontab格式或如`every 1 second|minute|hour|day|month|week`格式，更简单
* 执行时间进行了修正，会在秒/分开始的时候才执行，所以初次执行会有不到1秒/1分的延时
* 基于时间轮，保证实时性，任务容量会更高些。并发性能是github.com/robfig/cron的百倍以上，且更准时
* 回调函数可以增加一些参数，更容易调试，使用也更方便
* 支持整时执行和立即执行

## Install

```shell
go get -u github.com/lizongying/cron
```

## Usage

### job field

* Spec: 定时
* OnlyOnce: 只执行一次。默认false
* RunIfDelay: 即使超时(超过最大job处理数量)也会执行，否则本次不执行。默认false
* RunType: cron.Now 基于当前时间立即执行，默认; cron.Divisibility 整时运行。
* Id: 任务的唯一id，必须设置且不能重复。
* Meta: 任务的额外参数，非必须设置。
* Callback: 回调方法。

### cron options

* WithIntervalSecond 设置时间轮的间隔为秒，即定时任务最小间隔为一秒。此项为非默认设置。

```go
WithIntervalSecond() Options
```

* WithIntervalMinute 设置时间轮的间隔为分钟，即定时任务最小间隔为一分钟。此项为默认设置。

```go
WithIntervalMinute() Options
```

* WithLogger 设置使用自定义日志

```go
WithLogger(logger Logger) Options
```

* WithLoggerStdout 设置日志输出到控制台

```go
WithLoggerStdout() Options
```

### run

```go
package main

import (
	"github.com/lizongying/cron/cron"
	"time"
)

func main() {
	logger := cron.NewLoggerStdout()
	c := cron.New(cron.WithIntervalSecond(), cron.WithLoggerStdout())
	c.MustAddJob(&cron.Job{
		Spec: "every 3 seconds",
		Id:   1,
		Callback: func(id int, meta any) {
			logger.Info(id, meta, time.Now())
		},
	})
	logger.Info("now", time.Now())
	c.MustStart()

	select {}
}

```

### stop

```go
package main

import (
	"github.com/lizongying/cron/cron"
)

func main() {
	c := cron.New()
	c.MustStop()
}

```

## Tips

* 建议秒级别最大任务控制在1,000,000(Apple M1 Pro, 32 GB))以内，防止任务超时。可能支持更大数量，请自行测试。

## TODO
