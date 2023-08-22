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
	slotCount   int
	slots       []*map[uint32]*Job // [slot][jobId]job
	jobs        map[uint32]uint32  // [jobId]slot
	stopChannel chan struct{}
	running     bool
	locker      sync.Mutex
	logger      Logger
}

func New(options ...Options) (c *Cron) {
	c = &Cron{
		jobs:        make(map[uint32]uint32),
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
		c.logger = Logger(&LoggerNothing{})
	}

	c.slots = make([]*map[uint32]*Job, c.slotCount)

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
			c.logger.Info("stopped")
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

	c.running = true
	c.logger.Info("cron started")
	return
}

func (c *Cron) runJobs(jobs *map[uint32]*Job) {
	c.locker.Lock()
	defer c.locker.Unlock()

	for id, job := range *jobs {
		go func(job *Job) {
			defer func() {
				if err := recover(); err != nil {
					c.logger.Error("job run err:", err)
				}
			}()
			job.Callback()
		}(job)
		delete(*jobs, id)
		if job.OnlyOnce {
			delete(c.jobs, id)
			return
		}

		if err := c.saveJob(id, job); err != nil {
			c.logger.Error(err)
			continue
		}
		c.logger.Info("job next time:", id, job.nextTime)
	}
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

	id = c.id.Add(1)

	c.locker.Lock()
	defer c.locker.Unlock()

	if err = c.saveJob(id, job); err != nil {
		c.logger.Error(err)
		return
	}
	c.logger.Info("job next time:", id, job.nextTime)
	return
}

func (c *Cron) saveJob(id uint32, job *Job) (err error) {
	slot, err := job.Next(c.interval)
	if err != nil {
		return
	}

	if c.slots[slot] == nil {
		jobs := make(map[uint32]*Job)
		c.slots[slot] = &jobs
	}

	(*c.slots[slot])[id] = job
	c.jobs[id] = slot

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

	c.locker.Lock()
	defer c.locker.Unlock()

	slot, ok := c.jobs[id]
	if !ok {
		err = errors.New("job not exists")
		c.logger.Error(err)
		return
	}

	delete(c.jobs, id)
	delete(*c.slots[slot], id)

	c.logger.Info("job remove success")

	return
}
