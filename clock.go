package cron

import (
	"errors"
	"fmt"
	"time"
)

type Clock struct {
	year        uint16
	duration    time.Duration
	seconds     uint64
	secondFirst uint8
	secondLast  uint8
	second      uint8
	minutes     uint64
	minuteFirst uint8
	minuteLast  uint8
	minute      uint8
	hours       uint64
	hourFirst   uint8
	hourLast    uint8
	hour        uint8
	days        uint64
	dayFirst    uint8
	dayLast     uint8
	day         uint8
	months      uint64
	monthFirst  uint8
	monthLast   uint8
	month       uint8
	weeks       uint64
	weekFirst   uint8
	weekLast    uint8
	week        uint8
}

func NewClock(duration time.Duration, seconds uint64, minutes uint64, hours uint64, days uint64, months uint64, weeks uint64) (clock *Clock, err error) {
	clock = &Clock{
		duration: duration,
		seconds:  seconds,
		minutes:  minutes,
		hours:    hours,
		days:     days,
		months:   months,
		weeks:    weeks,
	}
	clock.secondFirst = clock.getFirst(clock.seconds)
	clock.minuteFirst = clock.getFirst(clock.minutes)
	clock.hourFirst = clock.getFirst(clock.hours)
	clock.dayFirst = clock.getFirst(clock.days)
	clock.monthFirst = clock.getFirst(clock.months)
	clock.weekFirst = clock.getFirst(clock.weeks)
	clock.secondLast = clock.getLast(clock.seconds)
	clock.minuteLast = clock.getLast(clock.minutes)
	clock.hourLast = clock.getLast(clock.hours)
	clock.dayLast = clock.getLast(clock.days)
	clock.monthLast = clock.getLast(clock.months)
	clock.weekLast = clock.getLast(clock.weeks)

	now := time.Now()
	clock.year = uint16(now.Year())
	clock.week = uint8(now.Weekday())
	clock.month = uint8(now.Month())
	clock.day = uint8(now.Day())
	clock.hour = uint8(now.Hour())
	clock.minute = uint8(now.Minute())
	clock.second = uint8(now.Second())

	clock.init("month")

	err = clock.initWithWeek()
	return
}

func (c *Clock) getLast(v uint64) (last uint8) {
	for i := 0; i < 60; i++ {
		if v&(1<<i) > 0 {
			last = uint8(i)
		}
	}
	return
}

func (c *Clock) getFirst(v uint64) (first uint8) {
	for i := 0; i < 60; i++ {
		if v&(1<<i) > 0 {
			first = uint8(i)
			break
		}
	}
	return
}

func (c *Clock) init(t string) {
	switch t {
	case "month":
		for i := 12; i > 0; i-- {
			if c.months&(1<<i) > 0 {
				v := uint8(i)
				if v <= c.month {
					c.month = v
					if v == c.month {
						c.init("day")
					} else {
						c.day = c.dayLast
						c.getWeek()
						c.hour = c.hourLast
						c.minute = c.minuteLast
						c.second = c.secondLast
					}
					return
				}
			}
		}
		// last year
		c.year--
		c.month = c.monthLast
		c.day = c.dayLast
		c.getWeek()
		c.hour = c.hourLast
		c.minute = c.minuteLast
		c.second = c.secondLast
	case "day":
		for i := 31; i > 0; i-- {
			if c.days&(1<<i) > 0 {
				v := uint8(i)
				if v <= c.day {
					c.day = v
					c.getWeek()
					if v == c.day {
						c.init("hour")
					} else {
						c.hour = c.hourLast
						c.minute = c.minuteLast
						c.second = c.secondLast
					}
					return
				}
			}
		}
		// last month
		c.day = c.dayLast
		c.getWeek()
		c.hour = c.hourLast
		c.minute = c.minuteLast
		c.second = c.secondLast
	case "hour":
		for i := 23; i >= 0; i-- {
			if c.hours&(1<<i) > 0 {
				v := uint8(i)
				if v <= c.hour {
					c.hour = v
					if v == c.hour {
						c.init("minute")
					} else {
						c.minute = c.minuteLast
						c.second = c.secondLast
					}
					return
				}
			}
		}
		// last day
		c.hour = c.hourLast
		c.minute = c.minuteLast
		c.second = c.secondLast
	case "minute":
		for i := 59; i >= 0; i-- {
			if c.minutes&(1<<i) > 0 {
				v := uint8(i)
				if v <= c.minute {
					c.minute = v
					if v == c.minute {
						c.init("second")
					} else {
						c.second = c.secondLast
					}
					return
				}
			}
		}
		// last hour
		c.minute = c.minuteLast
		c.second = c.secondLast
	case "second":
		for i := 59; i >= 0; i-- {
			if c.seconds&(1<<i) > 0 {
				v := uint8(i)
				if v <= c.second {
					c.second = v
					return
				}
			}
		}
		// last minute
		c.second = c.secondLast
	}

	return
}

func (c *Clock) initWithWeek() (err error) {
	year := c.year
	for ; ; c.getPrev("day") {
		// TODO better way
		if c.year-year > 1 {
			err = errors.New("no match")
			return
		}
		for i := 0; i < 7; i++ {
			if c.weeks&(1<<i) > 0 {
				v := uint8(i)
				if c.week == v {
					return
				}
			}
		}
		c.hour = c.hourFirst
		c.minute = c.minuteFirst
		c.second = c.secondFirst
	}
}

func (c *Clock) getPrev(t string) {
	switch t {
	case "month":
		for i := 12; i > 0; i-- {
			if c.months&(1<<i) > 0 {
				v := uint8(i)
				if v < c.month {
					c.month = v
					c.day = c.dayLast
					c.getWeek()
					c.hour = c.hourLast
					c.minute = c.minuteLast
					c.second = c.secondLast
					return
				}
			}
		}
		// prev year
		c.year--
		c.month = c.monthLast
		c.day = c.dayLast
		c.getWeek()
		c.hour = c.hourLast
		c.minute = c.minuteLast
		c.second = c.secondLast
	case "day":
		for i := 31; i > 0; i-- {
			if c.days&(1<<i) > 0 {
				v := uint8(i)
				if v < c.day {
					c.day = v
					c.getWeek()
					c.hour = c.hourLast
					c.minute = c.minuteLast
					c.second = c.secondLast
					return
				}
			}
		}
		// prev month
		c.getPrev("month")
	case "hour":
		for i := 23; i >= 0; i-- {
			if c.hours&(1<<i) > 0 {
				v := uint8(i)
				if v < c.hour {
					c.hour = v
					c.minute = c.minuteLast
					c.second = c.secondLast
					return
				}
			}
		}
		// prev day
		c.getPrev("day")
	case "minute":
		for i := 59; i >= 0; i-- {
			if c.minutes&(1<<i) > 0 {
				v := uint8(i)
				if v < c.minute {
					c.minute = v
					c.second = c.secondLast
					return
				}
			}
		}
		// prev hour
		c.getPrev("hour")
	case "second":
		for i := 59; i >= 0; i-- {
			if c.seconds&(1<<i) > 0 {
				v := uint8(i)
				if v < c.second {
					c.second = v
					return
				}
			}
		}
		// prev minute
		c.getPrev("minute")
	}

	return
}

func (c *Clock) getNext(t string) {
	switch t {
	case "month":
		for i := 1; i < 13; i++ {
			if c.months&(1<<i) > 0 {
				v := uint8(i)
				if v > c.month {
					c.month = v
					c.day = c.dayFirst
					c.getWeek()
					c.hour = c.hourFirst
					c.minute = c.minuteFirst
					c.second = c.secondFirst
					return
				}
			}
		}
		// next year
		c.year++
		c.month = c.monthFirst
		c.day = c.dayFirst
		c.getWeek()
		c.hour = c.hourFirst
		c.minute = c.minuteFirst
		c.second = c.secondFirst
	case "day":
		for i := 1; i < 32; i++ {
			if c.days&(1<<i) > 0 {
				v := uint8(i)
				if v > c.day {
					c.day = v
					c.getWeek()
					c.hour = c.hourFirst
					c.minute = c.minuteFirst
					c.second = c.secondFirst
					return
				}
			}
		}
		// next month
		c.getNext("month")
	case "hour":
		for i := 0; i < 24; i++ {
			if c.hours&(1<<i) > 0 {
				v := uint8(i)
				if v > c.hour {
					c.hour = v
					c.minute = c.minuteFirst
					c.second = c.secondFirst
					return
				}
			}
		}
		// next day
		c.getNext("day")
	case "minute":
		for i := 0; i < 60; i++ {
			if c.minutes&(1<<i) > 0 {
				v := uint8(i)
				if v > c.minute {
					c.minute = v
					c.second = c.secondFirst
					return
				}
			}
		}
		// next hour
		c.getNext("hour")
	case "second":
		for i := 0; i < 60; i++ {
			if c.seconds&(1<<i) > 0 {
				v := uint8(i)
				if v > c.second {
					c.second = v
					return
				}
			}
		}
		// next minute
		c.getNext("minute")
	}

	return
}

func (c *Clock) next() {
	if c.duration == time.Second {
		c.getNext("second")
	} else {
		c.getNext("minute")
	}

	return
}

func (c *Clock) NextWithWeek() (now time.Time, err error) {
	year := c.year
	c.next()

	for ; ; c.getNext("day") {
		// TODO better way
		if c.year-year > 2 {
			err = errors.New("no match")
			return
		}
		for i := 0; i < 7; i++ {
			if c.weeks&(1<<i) > 0 {
				v := uint8(i)
				if c.week == v {
					now = c.Now()
					return
				}
			}
		}
	}
}

func (c *Clock) getWeek() {
	day, _ := time.ParseInLocation(time.DateOnly, fmt.Sprintf("%d-%d-%d", c.year, c.month, c.day), time.Local)
	c.week = uint8(day.Weekday())
	return
}

func (c *Clock) Now() time.Time {
	now, _ := time.ParseInLocation(time.DateTime, fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", c.year, c.month, c.day, c.hour, c.minute, c.second), time.Local)
	return now
}

func (c *Clock) String() string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d %s", c.year, c.month, c.day, c.hour, c.minute, c.second, time.Weekday(c.week))
}
