# Cron

更简单的定时任务。

[cron-simple](https://github.com/lizongying/cron/tree/simple)
[cron-simple-v2](https://github.com/lizongying/cron/tree/simple_v2)

cron-simple: 降低内存使用率
cron-simple-v2: 降低cpu使用率

## Features

* 支持秒级定时，更准时。
* 通过执行周期、何时开始执行设置，更方便。
* 性能更高。
* 占用资源更少。

## Install

```shell
go get -u github.com/lizongying/cron@simple-v2
```

## Usage

### job method

* 设定执行周期，必须设置。

```go
EverySecond(v uint8) *Job
EveryMinute(v uint8) *Job
EveryHour(v uint8) *Job
EveryDay(v uint8) *Job
EveryMonth(v uint8) *Job
EveryWeek(v uint8) *Job
```

* 设定开始时间，非必须设置，如不设置，在cron开始后立即执行。

```go
// 忽略错误
MustSince(timeStr string) *Job

// 根据时间字符串，设定开始时间
// 比如定时04:05，每5分钟执行一次。
// 若当前时间03:06，会在04分05秒开始执行；
// 若当前时间05:02，会在09分05秒开始执行。
// 2006-01-02 15:04:05
// 01-02 15:04:05
// 02 15:04:05
// 15:04:05
// 04:05
// 05
// 未设置部分以当前时间填充
Since(timeStr string) (job *Job, err error)

// 直接设置开始时间
SinceTime(t time.Time) *Job

```

* 设定回调函数，必须设置。

```go
Callback(callback Callback) *Job
```

* 获取下次执行时间。

```go
NextTime() time.Time
```

### cron options

* 使用自定义日志

```go
WithLogger(logger Logger) Options
```

* 日志输出到控制台

```go
WithStdout() Options
```

### cron method

* 增加job

```go
MustAddJob(job *Job) (id uint32)
AddJob(job *Job) (id uint32, err error)
```

* 删除job

```go
MustRemoveJob(id uint32)
RemoveJob(id uint32) (err error)
```

* 修改job

```go
MustUpdateJob(id uint32, job *Job)
UpdateJob(id uint32, job *Job) (err error)
```

* 查询job

```go
MustGetJob(id uint32) (job *Job)
GetJob(id uint32) (job *Job, err error)
```

### run

```go
package main

import (
	"github.com/lizongying/cron"
)

func main() {
	logger := cron.NewLoggerStdout()
	c := cron.New(cron.WithStdout())
	c.MustStart()
	id := c.MustAddJob(new(cron.Job).
		EverySecond(10).
		MustSince("10:15").
		Callback(func() {
			logger.Info("callback")
		}))
	logger.Info("id", id)
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

* 建议秒级别最大任务控制在1,000,000(Apple M1 Pro, 32 GB))以内，防止任务超时。可能支持更大数量，请自行测试。

## Performance

结论：

和robfig/cron对比，相同数量任务，内存约为robfig/cron一半；

任务容量（任务不超时最大数量）约为robfig/cron的四倍

robfig/cron

```go
package main

import (
	cron "github.com/robfig/cron/v3"
	"log"
	"runtime"
	"time"
)

func main() {
	num := 1000000
	begin := time.Now()
	c := cron.New(cron.WithSeconds())
	c.Start()
	for i := 1; i <= num; i++ {
		v := i
		_, _ = c.AddFunc("@every 3s", func() {
			if v == num {
				var mem runtime.MemStats
				runtime.ReadMemStats(&mem)
				now := time.Now()
				log.Printf("Alloc = %v M, Spend = %v\n", mem.Alloc/1024/1024, now.Sub(begin))
				begin = now
			}
		})
	}
	log.Println("now", begin)
	select {}
}

```

lizongying/cron

```go
package main

import (
	"github.com/lizongying/cron"
	"log"
	"runtime"
	"time"
)

func main() {
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

```

lizongying 1,000,000:

![lizongying 1,000,000](./screenshot/lizongying_1000000.png)

robfig 1,000,000:

![robfig 1,000,000](./screenshot/robfig_1000000.png)

lizongying 4,000,000:

![lizongying 4,000,000](./screenshot/lizongying_4000000.png)

robfig 2,000,000:

如果任务执行时间一直超过3秒钟，可以认为到了最大容量。

robfig到了2,000,000就会出现任务超时情况。

![robfig 2,000,000](./screenshot/robfig_2000000.png)


