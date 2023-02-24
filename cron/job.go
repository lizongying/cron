package cron

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Callback func(id int, meta any)

type Job struct {
	Spec     string
	OnlyOnce bool
	Id       int
	Meta     any
	Callback Callback

	slot     int
	nextTime time.Time
	ok       bool
}

var reEvery = regexp.MustCompile(`@every\s(\d+)\s(second|minute|hour|day|month|week)s?`)

func (j *Job) Next(interval time.Duration) (err error) {
	now := time.Now()

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
			if !j.ok && v == 1 {
				now = now.Add(time.Second)
				j.ok = true
			}
			now = now.Add(time.Second * time.Duration(v))
		}
		if r[2] == "minute" {
			if v > 59 {
				err = errors.New("parse err")
				return
			}
			if !j.ok && v == 1 {
				now = now.Add(time.Minute)
				j.ok = true
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

	j.slot = GetSlot(now, interval)

	return
}

func GetSlot(t time.Time, interval time.Duration) (slot int) {
	if interval == time.Minute {
		slot |= t.Minute()
		slot |= t.Hour() << 6
		slot |= t.Day() << 11
		slot |= int(t.Month()) << 16
		slot |= int(t.Weekday()) << 20

		return
	}
	slot |= t.Second()
	slot |= t.Minute() << 6
	slot |= t.Hour() << 12
	slot |= t.Day() << 17
	slot |= int(t.Month()) << 22
	slot |= int(t.Weekday()) << 26

	return
}

func GetDateTime(slot int, interval time.Duration) (second int, minute int, hour int, day int, month int, week int) {
	if interval == time.Minute {
		minute = slot & 0x3f
		slot >>= 6
		hour = slot & 0x3f
		slot >>= 5
		day = slot & 0x1f
		slot >>= 5
		month = slot & 0x1f
		slot >>= 4
		week = slot & 0x7

		return
	}
	second = slot & 0x3f
	slot >>= 6
	minute = slot & 0x3f
	slot >>= 6
	hour = slot & 0x3f
	slot >>= 5
	day = slot & 0x1f
	slot >>= 5
	month = slot & 0x1f
	slot >>= 4
	week = slot & 0x7

	return
}
