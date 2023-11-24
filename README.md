# Cron

Simpler scheduling for tasks.

[cron](https://github.com/lizongying/cron/)
[cron-simple-v2](https://github.com/lizongying/cron/tree/simple_v2/)
[中文](./README_CN.md)

cron-simple-v2: Removed support for the crontab format, opting for a simpler approach. Additionally, reduced memory and CPU usage. It is recommended to use this branch.

## Features

* Supports second-level timing for more precise scheduling.
* Configurable through execution intervals and start times for added convenience.
* Improved performance.
* Minimal resource footprint.

## Install

```shell
go get -u github.com/lizongying/cron@simple-v2
```

## Usage

### job method

* Setting an execution interval is mandatory.

```go
// Ignoring errors.
MustEverySpec(spec string) *Job

// 1s/2i/3h/4d/5m/6w
EverySpec(spec string) error
EverySecond(v uint8) *Job
EveryMinute(v uint8) *Job
EveryHour(v uint8) *Job
EveryDay(v uint8) *Job
EveryMonth(v uint8) *Job
EveryWeek(v uint8) *Job
```

* Setting a start time is optional. If not set, the execution will take place immediately after the cron starts.

```go
// Ignoring errors.
MustSince(timeStr string) *Job

// Setting the start time based on a time string
// For example, scheduling at 04:05 with an interval of 5 minutes.
// If the current time is 03:06, execution will begin at 04:05.
// If the current time is 05:02, execution will begin at 09:05.
// Formats:
// - 2006-01-02 15:04:05
// - 01-02 15:04:05
// - 02 15:04:05
// - 15:04:05
// - 04:05
// - 05
// Unspecified parts are filled with the current time.
Since(timeStr string) error

// Directly setting the start time.
SinceTime(t time.Time) *Job

```

* Setting a callback function is mandatory.

```go
Callback(callback Callback) *Job
```

* Getting the next execution time.

```go
NextTime() time.Time
```

### cron options

* Logging output to the console.

```go
WithLogger(logger Logger) Options
```

* 日志输出到控制台

```go
WithStdout() Options
```

### cron method

* Adding a job.

```go
MustAddJob(job *Job) (id uint32)
AddJob(job *Job) (id uint32, err error)
```

* Deleting a job.

```go
MustRemoveJob(id uint32)
RemoveJob(id uint32) (err error)
```

* Modifying a job.

```go
MustUpdateJob(id uint32, job *Job)
UpdateJob(id uint32, job *Job) (err error)
```

* Querying a job.

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

* It is recommended to keep the maximum number of second-level tasks within 1,000,000 (for Apple M1 Pro, 32 GB) to
  prevent task timeouts. It might be possible to support a larger quantity, but it's advised to conduct your own
  testing.

## Performance

Conclusion:

Compared to robfig/cron, with the same number of tasks, the memory usage is approximately half.

The task capacity (maximum number of tasks without timeouts) is approximately four times that of robfig/cron.

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

If the task execution consistently exceeds 3 seconds, it can be considered reaching the maximum capacity.

For robfig/cron, task timeouts start occurring around 2,000,000 tasks.

![robfig 2,000,000](./screenshot/robfig_2000000.png)


