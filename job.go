package cron

import (
	"errors"
	"fmt"
	"time"
)

type Callback func()

type EveryType uint8

const (
	second EveryType = iota
	minute
	hour
	day
	month
	week
)

type Job struct {
	timestamp  uint32
	everyValue uint8
	everyType  EveryType
	callback   Callback
}

func (j *Job) NextTime() time.Time {
	return time.Unix(int64(j.timestamp), 0)
}
func (j *Job) EverySecond(v uint8) *Job {
	j.everyType = second
	j.everyValue = v
	return j
}
func (j *Job) EveryMinute(v uint8) *Job {
	j.everyType = minute
	j.everyValue = v
	return j
}
func (j *Job) EveryHour(v uint8) *Job {
	j.everyType = hour
	j.everyValue = v
	return j
}
func (j *Job) EveryDay(v uint8) *Job {
	j.everyType = day
	j.everyValue = v
	return j
}
func (j *Job) EveryMonth(v uint8) *Job {
	j.everyType = month
	j.everyValue = v
	return j
}
func (j *Job) EveryWeek(v uint8) *Job {
	j.everyType = week
	j.everyValue = v
	return j
}
func (j *Job) Callback(callback Callback) *Job {
	j.callback = callback
	return j
}
func (j *Job) MustSince(timeStr string) *Job {
	if _, err := j.Since(timeStr); err != nil {
		fmt.Println(err)
	}
	return j
}
func (j *Job) Since(timeStr string) (job *Job, err error) {
	l := len(timeStr)
	if l < 2 || l > 19 {
		err = errors.New("timeStr too short or too long")
		return
	}

	var t time.Time
	t, err = time.ParseInLocation(time.DateTime, time.Now().Format(time.DateTime)[:19-l]+timeStr, time.Local)
	if err != nil {
		return
	}

	j.SinceTime(t)
	job = j
	return
}
func (j *Job) SinceTime(t time.Time) *Job {
	j.timestamp = uint32(t.Unix())
	return j
}
func (j *Job) init(timestamp uint32) (err error) {
	if j == nil {
		err = errors.New("job nil")
		return
	}

	if j.everyValue == 0 {
		err = errors.New("everyValue 0")
		return
	}

	if j.callback == nil {
		err = errors.New("callback nil")
		return
	}

	if j.timestamp == 0 {
		j.timestamp = timestamp
	} else {
		for ; j.timestamp < timestamp; j.next() {
		}
	}
	return
}
func (j *Job) next() {
	now := time.Unix(int64(j.timestamp), 0)
	switch j.everyType {
	case second:
		j.timestamp = uint32(now.Add(time.Second * time.Duration(j.everyValue)).Unix())
	case minute:
		j.timestamp = uint32(now.Add(time.Minute * time.Duration(j.everyValue)).Unix())
	case hour:
		j.timestamp = uint32(now.Add(time.Hour * time.Duration(j.everyValue)).Unix())
	case day:
		j.timestamp = uint32(now.AddDate(0, 0, int(j.everyValue)).Unix())
	case month:
		j.timestamp = uint32(now.AddDate(0, int(j.everyValue), 0).Unix())
	case week:
		j.timestamp = uint32(now.AddDate(0, 0, 7*int(j.everyValue)).Unix())
	}
	return
}
