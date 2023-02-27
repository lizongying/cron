package cron

import (
	"errors"
	"sync"
	"time"
)

type status int

const (
	ready status = iota
	running
)

type cron struct {
	interval    time.Duration
	ticker      *time.Ticker
	slotCount   int
	slots       []*map[int]*Job // [slot][jobId]job
	jobs        map[int]int     // map[jobId]slot
	stopChannel chan struct{}
	status      status
	locker      sync.Mutex
	logger      Logger
}

func New(options ...Options) (c *cron) {
	c = &cron{
		jobs:        make(map[int]int),
		stopChannel: make(chan struct{}),
	}

	for _, v := range options {
		v(c)
	}

	if c.interval == 0 {
		c.interval = time.Minute
		c.slotCount = 60 * 24 * 366
	}
	if c.logger == nil {
		var logger Logger
		logger = &LoggerNothing{}
		c.logger = logger
	}

	c.slots = make([]*map[int]*Job, c.slotCount)

	return
}

func (c *cron) Start() (err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Println(err)
		return
	}

	if c.status == running {
		err = errors.New("cron already running")
		c.logger.Println(err)
		return
	}

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

	go func() {
		defer func() {
			c.logger.Println("stopped")
		}()

		slot := GetSlotSinceYear(now, c.interval)
		jobs := c.slots[slot]
		if jobs != nil && len(*jobs) > 0 {
			go c.runJobs(jobs)
		}

		for {
			select {
			case now = <-c.ticker.C:
				slot = GetSlotSinceYear(now, c.interval)
				jobs = c.slots[slot]
				if jobs != nil && len(*jobs) > 0 {
					go c.runJobs(jobs)
				}
			case <-c.stopChannel:
				return
			}
		}
	}()

	c.status = running
	c.logger.Println("cron started")
	return
}

func (c *cron) runJobs(jobs *map[int]*Job) {
	for _, job := range *jobs {
		go job.Callback(job.Id, job.Meta, job.nextTime)
		delete(*jobs, job.Id)

		if job.OnlyOnce {
			delete(c.jobs, job.Id)
			continue
		}
		err := c.saveJob(job)
		if err != nil {
			c.logger.Println(err)
			continue
		}
		c.logger.Println("job next time:", job.Id, job.nextTime)
	}
}

func (c *cron) Stop() (err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Println(err)
		return
	}

	if c.status != running {
		err = errors.New("cron not running")
		c.logger.Println(err)
		return
	}

	close(c.stopChannel)
	c.ticker.Stop()

	c.status = ready
	c.logger.Println("cron stopped")

	return
}

func (c *cron) AddJob(job *Job) (err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Println(err)
		return
	}

	if job.Spec == "" {
		err = errors.New("spect empty")
		c.logger.Println(err)
		return
	}

	if job.Id == 0 {
		err = errors.New("id is 0")
		c.logger.Println(err)
		return
	}

	if job.Callback == nil {
		err = errors.New("callback is nil")
		c.logger.Println(err)
		return
	}

	now := time.Now()
	if c.interval == time.Second {
		now = time.Unix(now.Unix(), 0)
	} else {
		now = time.Unix(now.Unix()-int64(now.Second()), 0)
	}
	job.nextTime = now

	if err = c.saveJob(job); err != nil {
		c.logger.Println(err)
		return
	}

	return
}

func (c *cron) saveJob(job *Job) (err error) {
	prevTime := job.nextTime
	slot, err := job.Next(c.interval)
	if err != nil {
		err = errors.New("job parse err")
		c.logger.Println(err)
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	if c.slots[slot] == nil {
		jobs := make(map[int]*Job)
		c.slots[slot] = &jobs
	}

	(*c.slots[slot])[job.Id] = job
	c.jobs[job.Id] = slot

	c.logger.Println("job save success:", job.Id, prevTime)

	return
}

func (c *cron) RemoveJob(id int) (err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Println(err)
		return
	}

	if id == 0 {
		err = errors.New("id 0")
		c.logger.Println(err)
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	slot, ok := c.jobs[id]
	if !ok {
		err = errors.New("job not exists")
		c.logger.Println(err)
		return
	}

	delete(c.jobs, id)
	delete(*c.slots[slot], id)

	c.logger.Println("job remove success")

	return
}
