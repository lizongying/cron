package cron

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Save struct {
	Id  uint32
	Job *Job
}

type Cron struct {
	running     bool
	id          atomic.Uint32
	locker      sync.RWMutex
	jobs        map[uint32]*Job
	stopChannel chan struct{}
	logger      Logger

	removeChannel chan uint32
	saveChannel   chan Save
}

func New(options ...Options) (c *Cron) {
	c = &Cron{
		jobs:        make(map[uint32]*Job),
		stopChannel: make(chan struct{}),
		logger:      Logger(new(LoggerNothing)),

		removeChannel: make(chan uint32),
		saveChannel:   make(chan Save),
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

	timer := time.NewTimer(time.Hour * 24 * 365)
	jobs := make(map[uint32]*Job)
	var now time.Time
	var nextTime time.Time

	timestamp := uint32(time.Now().Unix() + 1)
	for _, job := range jobs {
		if err = job.init(timestamp); err != nil {
			c.logger.Error(err)
			continue
		}
	}
	for id, job := range c.jobs {
		if err = job.init(timestamp); err != nil {
			c.logger.Error(err)
			continue
		}

		if len(jobs) != 0 {
			timestamp = job.timestamp
			jobs[id] = job
			continue
		}
		if job.timestamp < timestamp {
			timestamp = job.timestamp
			jobs = make(map[uint32]*Job)
			jobs[id] = job
			continue
		}
		if job.timestamp == timestamp {
			jobs[id] = job
		}
	}

	if len(jobs) != 0 {
		now = time.Now()
		nextTime = time.Unix(int64(timestamp), 0)
		if nextTime.After(now) {
			timer.Reset(nextTime.Sub(now))
		} else {
			timer.Reset(0)
		}
	}

	go func() {
		for {
			select {
			case <-timer.C:
				for _, job := range jobs {
					go c.runJob(job)
				}
				for _, job := range jobs {
					job.next()
				}

				jobs = make(map[uint32]*Job)
				for id, job := range c.jobs {
					if len(jobs) == 0 {
						timestamp = job.timestamp
						jobs[id] = job
						continue
					}
					if job.timestamp < timestamp {
						timestamp = job.timestamp
						jobs = make(map[uint32]*Job)
						jobs[id] = job
						continue
					}
					if job.timestamp == timestamp {
						jobs[id] = job
					}
				}

				if len(jobs) == 0 {
					continue
				}

				now = time.Now()
				nextTime = time.Unix(int64(timestamp), 0)
				if nextTime.After(now) {
					timer.Reset(nextTime.Sub(now))
				} else {
					timer.Reset(0)
				}
			case save := <-c.saveChannel:
				id := save.Id
				job := save.Job
				if err = job.init(timestamp); err != nil {
					c.logger.Error(err)
					continue
				}
				c.locker.Lock()
				c.jobs[id] = job
				c.locker.Unlock()

				if len(jobs) == 0 {
					timestamp = job.timestamp
					jobs[id] = job
				}
				if job.timestamp < timestamp {
					timestamp = job.timestamp
					jobs = make(map[uint32]*Job)
					jobs[id] = job
				}
				if job.timestamp == timestamp {
					jobs[id] = job
				}

				if len(jobs) == 0 {
					continue
				}

				now = time.Now()
				nextTime = time.Unix(int64(timestamp), 0)
				if nextTime.After(now) {
					timer.Reset(nextTime.Sub(now))
				} else {
					timer.Reset(0)
				}
			case id := <-c.removeChannel:
				c.locker.Lock()
				delete(c.jobs, id)
				c.locker.Unlock()

				delete(jobs, id)

				if len(jobs) == 0 {
					for id, job := range c.jobs {
						if len(jobs) == 0 {
							timestamp = job.timestamp
							jobs[id] = job
							continue
						}
						if job.timestamp < timestamp {
							timestamp = job.timestamp
							jobs = make(map[uint32]*Job)
							jobs[id] = job
							continue
						}
						if job.timestamp == timestamp {
							jobs[id] = job
						}
					}

					if len(jobs) == 0 {
						continue
					}

					now = time.Now()
					nextTime = time.Unix(int64(timestamp), 0)
					if nextTime.After(now) {
						timer.Reset(nextTime.Sub(now))
					} else {
						timer.Reset(0)
					}
				}
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

	c.id.Add(1)
	id = c.id.Load()
	c.saveChannel <- Save{
		Id:  id,
		Job: job,
	}
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

	if id == 0 {
		err = errors.New("id zero")
		c.logger.Error(err)
		return
	}

	c.removeChannel <- id
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

	c.saveChannel <- Save{
		Id:  id,
		Job: job,
	}
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
