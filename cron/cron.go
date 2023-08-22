package cron

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Cron struct {
	id          atomic.Uint32
	interval    time.Duration
	ticker      *time.Ticker
	jobs        []*Job
	stopChannel chan struct{}
	running     bool
	locker      sync.Mutex
	logger      Logger
	lastMoment  uint32
}

func New(options ...Options) (c *Cron) {
	c = &Cron{
		stopChannel: make(chan struct{}),
	}

	for _, v := range options {
		v(c)
	}

	if c.interval == 0 {
		c.interval = time.Minute
	}

	if c.logger == nil {
		c.logger = Logger(&LoggerNothing{})
	}

	return
}

func (c *Cron) MustStart() {
	if err := c.Start(); err != nil {
		c.logger.Error(err)
	}
}

func (c *Cron) Start() (err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Error(err)
		return
	}

	if c.running {
		err = errors.New("cron already running")
		c.logger.Error(err)
		return
	}

	go func() {
		defer func() {
			c.logger.Info("stopped")
		}()

		c.lastMoment = LastMoment(c.interval)

		now := time.Now()
		var nextTime time.Time
		if c.interval == time.Second {
			nextTime = time.Unix(now.Unix(), 0).Add(time.Second)
			time.Sleep(time.Duration(nextTime.UnixNano() - now.UnixNano()))
		} else {
			nextTime = time.Unix(now.Unix()-int64(now.Second()), 0).Add(time.Minute)
			time.Sleep(time.Duration(nextTime.UnixNano() - now.UnixNano()))
		}
		now = nextTime
		c.ticker = time.NewTicker(c.interval)

		slot := SlotSinceYear(now, c.interval)
		for _, job := range c.jobs {
			if job.Deleted {
				continue
			}
			if err = job.Next(c.interval); err != nil {
				c.logger.Error(err)
				continue
			}
			if job.Slot() == slot {
				go c.runJob(job)
			}
		}

		for {
			select {
			case now = <-c.ticker.C:
				if slot == c.lastMoment {
					slot = 0
					c.lastMoment = LastMoment(c.interval)
				} else {
					slot++
				}
				for _, job := range c.jobs {
					if job.Deleted {
						continue
					}
					if job.Slot() == slot {
						go c.runJob(job)
					}
				}
			case <-c.stopChannel:
				return
			}
		}
	}()

	c.running = true
	c.logger.Info("cron started")
	return
}

func (c *Cron) runJob(job *Job) {
	defer func() {
		if err := recover(); err != nil {
			c.logger.Error("job run err:", err)
		}
	}()
	go func() {
		err := job.Next(c.interval)
		if err != nil {
			c.logger.Error(err)
		}
	}()
	job.Callback()
}

func (c *Cron) MustStop() {
	if err := c.Stop(); err != nil {
		c.logger.Error(err)
	}
}

func (c *Cron) Stop() (err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Error(err)
		return
	}

	if !c.running {
		err = errors.New("cron not running")
		c.logger.Error(err)
		return
	}

	close(c.stopChannel)
	c.ticker.Stop()

	c.running = false
	c.logger.Info("cron stopped")

	return
}

func (c *Cron) MustAddJob(spec string, job *Job) (id uint32) {
	var err error
	id, err = c.AddJob(spec, job)
	if err != nil {
		c.logger.Error(err)
	}
	return
}

func (c *Cron) AddJob(spec string, job *Job) (id uint32, err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Error(err)
		return
	}

	if spec == "" {
		err = errors.New("spec empty")
		c.logger.Error(err)
		return
	}

	if job.Callback == nil {
		err = errors.New("callback is nil")
		c.logger.Error(err)
		return
	}

	if err = job.Init(spec, c.interval); err != nil {
		c.logger.Error(err)
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	c.jobs = append(c.jobs, job)

	c.logger.Info("job next time:", c.id.Load(), job.nextTime)
	c.id.Add(1)
	return
}

func (c *Cron) MustRemoveJob(id uint32) {
	if err := c.RemoveJob(id); err != nil {
		c.logger.Error(err)
	}
}

func (c *Cron) RemoveJob(id uint32) (err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Error(err)
		return
	}

	if id == 0 {
		err = errors.New("id 0")
		c.logger.Error(err)
		return
	}

	if id > c.id.Load() {
		err = errors.New("job not exists")
		c.logger.Error(err)
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	c.jobs[id].Deleted = true
	c.logger.Info("job remove success")
	return
}
