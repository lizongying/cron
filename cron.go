package cron

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Cron struct {
	id          atomic.Uint32
	running     bool
	jobs        map[uint32]*Job
	stopChannel chan struct{}
	logger      Logger
	locker      sync.RWMutex
}

func New(options ...Options) (c *Cron) {
	c = &Cron{
		jobs:        make(map[uint32]*Job),
		stopChannel: make(chan struct{}),
		logger:      Logger(new(LoggerNothing)),
	}

	for _, v := range options {
		v(c)
	}

	return
}

func (c *Cron) MustStart() {
	_ = c.Start()
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

	now := time.Unix(time.Now().Unix(), 0).Add(time.Second)
	timer := time.NewTimer(now.Sub(time.Now()))

	go func() {
		timestamp := uint32(now.Unix())
		for _, job := range c.jobs {
			if err = job.init(timestamp); err != nil {
				c.logger.Error(err)
				continue
			}
		}

		for {
			select {
			case <-timer.C:
				for _, job := range c.jobs {
					if job.timestamp == timestamp {
						go c.runJob(job)
					}
				}
				timestamp++
				now = now.Add(time.Second)
				timer.Reset(now.Sub(time.Now()))
			case <-c.stopChannel:
				timer.Stop()
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

	go job.next()

	job.callback()
}

func (c *Cron) MustStop() {
	_ = c.Stop()
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

	c.running = false
	c.logger.Info("cron stopped")
	return
}

func (c *Cron) MustAddJob(job *Job) (id uint32) {
	id, _ = c.AddJob(job)
	return
}

func (c *Cron) AddJob(job *Job) (id uint32, err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Error(err)
		return
	}

	if err = job.init(uint32(time.Now().Unix())); err != nil {
		c.logger.Error(err)
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	c.jobs[c.id.Load()] = job
	c.id.Add(1)
	return
}

func (c *Cron) MustRemoveJob(id uint32) {
	_ = c.RemoveJob(id)
}

func (c *Cron) RemoveJob(id uint32) (err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Error(err)
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	delete(c.jobs, id)
	c.logger.Info("job remove success")
	return
}

func (c *Cron) MustUpdateJob(id uint32, job *Job) {
	_ = c.UpdateJob(id, job)
}

func (c *Cron) UpdateJob(id uint32, job *Job) (err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Error(err)
		return
	}

	if _, err = c.GetJob(id); err != nil {
		err = errors.New("job not exists")
		c.logger.Error(err)
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	if err = job.init(uint32(time.Now().Unix())); err != nil {
		c.logger.Error(err)
		return
	}
	c.jobs[id] = job
	c.logger.Info("job update success")
	return
}

func (c *Cron) MustGetJob(id uint32) (job *Job) {
	job, _ = c.GetJob(id)
	return
}

func (c *Cron) GetJob(id uint32) (job *Job, err error) {
	if c == nil {
		err = errors.New("cron nil")
		c.logger.Error(err)
		return
	}

	c.locker.RLock()
	defer c.locker.RUnlock()

	var ok bool
	job, ok = c.jobs[id]
	if !ok {
		err = errors.New("job not exists")
		c.logger.Error(err)
		return
	}

	return
}
