package cron

import (
	"errors"
	"sync"
	"time"
)

type Status int

const (
	Ready Status = iota
	Running
)

type cron struct {
	interval    time.Duration
	ticker      *time.Ticker
	slotCount   int
	slots       []*map[int]*Job // [slot][jobId]job
	jobs        map[int]int     // map[jobId]slot
	slot        int             // current slot
	stopChannel chan struct{}
	status      Status
	locker      sync.Mutex
	logger      Logger
}

func NewCron(options ...Options) (c *cron) {
	c = &cron{
		jobs:        make(map[int]int),
		stopChannel: make(chan struct{}),
	}

	for _, v := range options {
		v(c)
	}

	if c.interval == 0 {
		c.interval = time.Minute
		c.slotCount = 1 << 23
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

	if c.status == Running {
		err = errors.New("cron already running")
		c.logger.Println(err)
		return
	}

	now := time.Now()
	if c.interval == time.Second {
		time.Sleep(time.Duration(time.Unix(now.Unix(), 0).Add(time.Second).UnixNano() - now.UnixNano()))
	} else {
		time.Sleep(time.Duration(time.Unix(now.Unix()-int64(now.Second()), 0).Add(time.Minute).UnixNano() - now.UnixNano()))
	}

	c.ticker = time.NewTicker(c.interval)

	go func() {
		defer func() {
			c.logger.Println("stopped")
		}()

		for {
			select {
			case now = <-c.ticker.C:
				if c.interval == time.Minute {
					now = time.Unix(now.Unix()-int64(now.Second()), 0)
				}
				slot := GetSlot(now, c.interval)
				jobs := c.slots[slot]
				if jobs != nil && len(*jobs) > 0 {
					go c.runJobs(jobs)
				}
			case <-c.stopChannel:
				return
			}
		}
	}()

	c.status = Running
	c.logger.Println("cron started")
	return
}

func (c *cron) runJobs(jobs *map[int]*Job) {
	for _, job := range *jobs {
		go job.Callback(job.Id, job.Meta)
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

	if c.status != Running {
		err = errors.New("cron not running")
		c.logger.Println(err)
		return
	}

	close(c.stopChannel)
	c.ticker.Stop()

	c.status = Ready
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

	if err = c.saveJob(job); err != nil {
		c.logger.Println(err)
		return
	}

	return
}

func (c *cron) saveJob(job *Job) (err error) {
	if err = job.Next(c.interval); err != nil {
		err = errors.New("job parse err")
		c.logger.Println(err)
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	if c.slots[job.slot] == nil {
		jobs := make(map[int]*Job)
		c.slots[job.slot] = &jobs
	}

	(*c.slots[job.slot])[job.Id] = job
	c.jobs[job.Id] = job.slot

	c.logger.Println("job save success:", job.Id)

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
