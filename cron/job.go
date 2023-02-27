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

type Job struct {
	Spec     string
	OnlyOnce bool
	Id       int
	Meta     any
	Callback Callback

	nextTime time.Time
}

var reEvery = regexp.MustCompile(`@every\s(\d+)\s(second|minute|hour|day|month|week)s?`)

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

func GetSlot(now time.Time, interval time.Duration) (slot int) {
	if interval == time.Minute {
		slot |= now.Minute()
		slot |= now.Hour() << 6
		slot |= now.Day() << 11
		slot |= int(now.Month()) << 16
		slot |= int(now.Weekday()) << 20

		return
	}
	slot |= now.Second()
	slot |= now.Minute() << 6
	slot |= now.Hour() << 12
	slot |= now.Day() << 17
	slot |= int(now.Month()) << 22
	slot |= int(now.Weekday()) << 26

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

func GetSlotSinceYear(now time.Time, interval time.Duration) (slot int) {
	year, _ := time.ParseInLocation("2006", now.Format("2006"), time.Local)
	if interval == time.Minute {
		slot = int(math.Floor(now.Sub(year).Minutes()))

		return
	}
	slot = int(math.Floor(now.Sub(year).Seconds()))

	return
}
