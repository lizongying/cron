# Cron

基于时间轮实现的定时任务，更准时，并发性能更高。支持crontab格式或`every 1 second|minute|hour|day|month|week`格式

[cron](https://github.com/lizongying/cron)

如果没有crontab格式要求，建议使用simple版本，占用内存更少。
[cron-simple-v2](https://github.com/lizongying/cron/tree/simple-v2)

## Features

* 基于时间轮实现，更准时，并发性能更高。
* 支持crontab或`every 1 second|minute|hour|day|month|week`格式
* 修正执行时间，会在整秒/分开始的时候才执行，所以初次执行会有不到1秒/1分的延时
* 支持立即或整时执行

## Install

```shell
go get -u github.com/lizongying/cron
```

## Usage

### job field

* Divisibility: 整时执行，默认false。
* Callback: 回调方法。

### cron options

* WithSecond 设置时间轮的间隔为秒，即定时任务最小间隔为一秒。此项为非默认设置。

```go
WithSecond() Options
```

* WithMinute 设置时间轮的间隔为分钟，即定时任务最小间隔为一分钟。此项为默认设置。

```go
WithMinute() Options
```

* WithLogger 设置使用自定义日志

```go
WithLogger(logger Logger) Options
```

* WithStdout 设置日志输出到控制台

```go
WithStdout() Options
```

* WithDivisibility 设置整点执行

```go
WithDivisibility() Options

```

### run

```go
package main

import (
	"github.com/lizongying/cron"
	"time"
)

func main() {
	begin := time.Now()
	logger := cron.NewLoggerStdout()
	c := cron.New(cron.WithSecond(), cron.WithStdout())
	id := c.MustAddJob("every 3 seconds", func() {
		logger.Info(time.Now().Sub(begin))
	})
	logger.Info("id", id)
	c.MustStart()
	select {}
}

```

### stop

```go
package main

import (
	"github.com/lizongying/cron"
)

func main() {
	c := cron.New()
	c.MustStop()
}

```

## Tips

* 建议秒级别最大任务控制在4,000,000(Apple M1 Pro, 32 GB))以内，防止任务超时。可能支持更大数量，请自行测试。
* 经测试，内存和并发性能均优于github.com/robfig/cron，可参考[cron-simple](https://github.com/lizongying/cron/tree/simple)

## TODO
