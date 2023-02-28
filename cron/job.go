package cron

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Callback func(id int, meta any, now time.Time)

type RunType int

const (
	Now RunType = iota
	Divisibility
)

type Job struct {
	Spec     string
	OnlyOnce bool
	RunType  RunType
	Id       int
	Meta     any
	Callback Callback

	nextTime time.Time
}

var reEvery = regexp.MustCompile(`@every\s(\d+)\s(second|minute|hour|day|month|week)s?`)

func (j *Job) Init(now time.Time, interval time.Duration) (err error) {
	if interval == time.Second {
		now = time.Unix(now.Unix(), 0)
	} else {
		now = time.Unix(now.Unix()-int64(now.Second()), 0)
	}
	r := reEvery.FindStringSubmatch(j.Spec)
	if len(r) == 3 {
		if strings.Index("second|minute|hour|day|month|week", r[2]) < 0 {
			err = errors.New("parse err")
			return
		}

		v, _ := strconv.Atoi(r[1])

		if r[2] == "second" {
			if v > 59 {
				err = errors.New("parse err")
				return
			}
			if j.RunType == Divisibility {
				now = now.Add(-time.Second * time.Duration(now.Second()%v))
			}
		}
		if r[2] == "minute" {
			if v > 59 {
				err = errors.New("parse err")
				return
			}
			if j.RunType == Divisibility {
				now = now.Add(-time.Second * time.Duration(now.Second()))
				now = now.Add(-time.Minute * time.Duration(now.Minute()%v))
			}
		}
		if r[2] == "hour" {
			if v > 23 {
				err = errors.New("parse err")
				return
			}
			if j.RunType == Divisibility {
				now = now.Add(-time.Second * time.Duration(now.Second()))
				now = now.Add(-time.Minute * time.Duration(now.Minute()))
				now = now.Add(-time.Hour * time.Duration(now.Hour()%v))
			}
		}
		if r[2] == "day" {
			if v > 30 {
				err = errors.New("parse err")
				return
			}
			if j.RunType == Divisibility {
				now = now.Add(-time.Second * time.Duration(now.Second()))
				now = now.Add(-time.Minute * time.Duration(now.Minute()))
				now = now.Add(-time.Hour * time.Duration(now.Hour()))
				now = now.AddDate(0, 0, -now.Day()%v)
			}
		}
		if r[2] == "month" {
			if v > 11 {
				err = errors.New("parse err")
				return
			}
			if j.RunType == Divisibility {
				now = now.Add(-time.Second * time.Duration(now.Second()))
				now = now.Add(-time.Minute * time.Duration(now.Minute()))
				now = now.Add(-time.Hour * time.Duration(now.Hour()))
				now = now.AddDate(0, 0, -(now.Day() - 1))
				now = now.AddDate(0, -int(now.Month())%v, 0)
			}
		}
		if r[2] == "week" {
			if v > 3 {
				err = errors.New("parse err")
				return
			}

			// default run on sunday
			// monday now = now.AddDate(0, 0, -int(now.Weekday())%7+1)
			if j.RunType == Divisibility {
				now = now.Add(-time.Second * time.Duration(now.Second()))
				now = now.Add(-time.Minute * time.Duration(now.Minute()))
				now = now.Add(-time.Hour * time.Duration(now.Hour()))
				now = now.AddDate(0, 0, -int(now.Weekday())%7)
			}
		}
		if interval == time.Minute {
			now = time.Unix(now.Unix()-int64(now.Second()), 0)
		}
	} else {
		err = errors.New("parse err")
		return
	}

	j.nextTime = now

	return
}

func (j *Job) Next(interval time.Duration) (slot int, err error) {
	now := j.nextTime
	r := reEvery.FindStringSubmatch(j.Spec)
	if len(r) == 3 {
		if strings.Index("second|minute|hour|day|month|week", r[2]) < 0 {
			err = errors.New("parse err")
			return
		}

		v, _ := strconv.Atoi(r[1])

		if r[2] == "second" {
			if v > 59 {
				err = errors.New("parse err")
				return
			}
			now = now.Add(time.Second * time.Duration(v))
		}
		if r[2] == "minute" {
			if v > 59 {
				err = errors.New("parse err")
				return
			}
			now = now.Add(time.Minute * time.Duration(v))
		}
		if r[2] == "hour" {
			if v > 23 {
				err = errors.New("parse err")
				return
			}
			now = now.Add(time.Hour * time.Duration(v))
		}
		if r[2] == "day" {
			if v > 30 {
				err = errors.New("parse err")
				return
			}
			now = now.AddDate(0, 0, v)
		}
		if r[2] == "month" {
			if v > 11 {
				err = errors.New("parse err")
				return
			}
			now = now.AddDate(0, v, 0)
		}
		if r[2] == "week" {
			if v > 3 {
				err = errors.New("parse err")
				return
			}
			now = now.AddDate(0, 0, 7*v)
		}
		if interval == time.Minute {
			now = time.Unix(now.Unix()-int64(now.Second()), 0)
		}
	} else {
		err = errors.New("parse err")
		return
	}

	j.nextTime = now

	slot = GetSlotSinceYear(now, interval)

	return
}

func GetSlotSinceYear(now time.Time, interval time.Duration) (slot int) {
	year, _ := time.ParseInLocation("2006", now.Format("2006"), time.Local)
	if interval == time.Minute {
		slot = int(math.Floor(now.Sub(year).Minutes()))

		return
	}
	slot = int(math.Floor(now.Sub(year).Seconds()))

	return
}
