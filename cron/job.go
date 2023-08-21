package cron

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Callback func(id int, meta any)

type RunType int

const (
	Now RunType = iota
	Divisibility
)

var reEvery = regexp.MustCompile(`every\s(\d+)\s(second|minute|hour|day|month|week)s?`)
var reDash = regexp.MustCompile(`(\d+)-(\d+)`)
var reSlash = regexp.MustCompile(`\*/(\d+)`)

type element struct {
	min  int
	max  int
	name string
}

var parser = []element{
	{0, 59, "second"},
	{0, 59, "minute"},
	{0, 23, "hour"},
	{1, 31, "day"},
	{1, 12, "month"},
	{0, 6, "week"},
}

type Job struct {
	Spec       string
	OnlyOnce   bool
	RunIfDelay bool
	RunType    RunType
	Id         int
	Meta       any
	Callback   Callback

	nextTime   time.Time
	clock      *Clock
	everyName  string
	everyValue int
}

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

		v, e := strconv.Atoi(r[1])
		if e != nil {
			err = errors.New("parse err")
			return
		}

		if r[2] == "second" {
			if v > 59 {
				err = errors.New("parse err")
				return
			}
			if j.RunType == Divisibility {
				now = now.Add(-time.Second * time.Duration(now.Second()%v))
			}
		} else if r[2] == "minute" {
			if v > 59 {
				err = errors.New("parse err")
				return
			}
			if j.RunType == Divisibility {
				now = now.Add(-time.Second * time.Duration(now.Second()))
				now = now.Add(-time.Minute * time.Duration(now.Minute()%v))
			}
		} else if r[2] == "hour" {
			if v > 23 {
				err = errors.New("parse err")
				return
			}
			if j.RunType == Divisibility {
				now = now.Add(-time.Second * time.Duration(now.Second()))
				now = now.Add(-time.Minute * time.Duration(now.Minute()))
				now = now.Add(-time.Hour * time.Duration(now.Hour()%v))
			}
		} else if r[2] == "day" {
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
		} else if r[2] == "month" {
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
		} else if r[2] == "week" {
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
		} else {
			err = errors.New("parse err")
			return
		}
		j.everyName = r[2]
		j.everyValue = v
	} else {
		li := strings.Split(j.Spec, " ")
		if len(li) == 5 {
			li = append([]string{"*"}, li...)
		}
		if len(li) == 6 {
			var list [6]uint64
		LOOP1:
			for i, v := range li {
				if v == "*" {
					begin := parser[i].min
					end := parser[i].max + 1
					for ii := begin; ii < end; ii++ {
						list[i] |= 1 << ii
					}
					continue
				}
				r = reSlash.FindStringSubmatch(v)
				if len(r) == 2 {
					every, e := strconv.Atoi(r[1])
					if e != nil {
						err = errors.New("parse err")
						break LOOP1
					}
					// TODO max?
					if every < 1 {
						err = errors.New("parse err")
						break LOOP1
					}

					if i == 0 {
						if every > 59 {
							err = errors.New("parse err")
							break LOOP1
						}
						if j.RunType == Divisibility {
							now = now.Add(-time.Second * time.Duration(now.Second()%every))
						}
					} else if i == 1 {
						if every > 59 {
							err = errors.New("parse err")
							break LOOP1
						}
						if j.RunType == Divisibility {
							now = now.Add(-time.Second * time.Duration(now.Second()))
							now = now.Add(-time.Minute * time.Duration(now.Minute()%every))
						}
					} else if i == 2 {
						if every > 23 {
							err = errors.New("parse err")
							break LOOP1
						}
						if j.RunType == Divisibility {
							now = now.Add(-time.Second * time.Duration(now.Second()))
							now = now.Add(-time.Minute * time.Duration(now.Minute()))
							now = now.Add(-time.Hour * time.Duration(now.Hour()%every))
						}
					} else if i == 3 {
						if every > 30 {
							err = errors.New("parse err")
							break LOOP1
						}
						if j.RunType == Divisibility {
							now = now.Add(-time.Second * time.Duration(now.Second()))
							now = now.Add(-time.Minute * time.Duration(now.Minute()))
							now = now.Add(-time.Hour * time.Duration(now.Hour()))
							now = now.AddDate(0, 0, -now.Day()%every)
						}
					} else if i == 4 {
						if every > 11 {
							err = errors.New("parse err")
							break LOOP1
						}
						if j.RunType == Divisibility {
							now = now.Add(-time.Second * time.Duration(now.Second()))
							now = now.Add(-time.Minute * time.Duration(now.Minute()))
							now = now.Add(-time.Hour * time.Duration(now.Hour()))
							now = now.AddDate(0, 0, -(now.Day() - 1))
							now = now.AddDate(0, -int(now.Month())%every, 0)
						}
					} else if i == 5 {
						if every > 3 {
							err = errors.New("parse err")
							break LOOP1
						}

						// default run on sunday
						// monday now = now.AddDate(0, 0, -int(now.Weekday())%7+1)
						if j.RunType == Divisibility {
							now = now.Add(-time.Second * time.Duration(now.Second()))
							now = now.Add(-time.Minute * time.Duration(now.Minute()))
							now = now.Add(-time.Hour * time.Duration(now.Hour()))
							now = now.AddDate(0, 0, -int(now.Weekday())%7)
						}
					} else {
						err = errors.New("parse err")
						break LOOP1
					}

					begin := parser[i].min
					end := parser[i].max + 1
					for ii := begin; ii < end; ii++ {
						if j.RunType == Divisibility {
							if ii%every == 0 {
								list[i] |= 1 << ii
							}
						} else {
							// TODO
							if ii%every == 0 {
								list[i] |= 1 << ii
							}
						}
					}
					continue
				}

				li2 := strings.Split(v, ",")
				for _, v2 := range li2 {
					r = reDash.FindStringSubmatch(v2)
					if len(r) == 3 {
						begin, e := strconv.Atoi(r[1])
						if e != nil {
							err = errors.New("parse err")
							break LOOP1
						}
						if begin > parser[i].max || begin < parser[i].min {
							err = errors.New("parse err")
							break LOOP1
						}
						end, e := strconv.Atoi(r[2])
						if e != nil {
							err = errors.New("parse err")
							break LOOP1
						}
						if end > parser[i].max || end < parser[i].min {
							err = errors.New("parse err")
							break LOOP1
						}

						end++
						for ii := begin; ii < end; ii++ {
							list[i] |= 1 << ii
						}

						continue
					}

					ii, e := strconv.Atoi(v2)
					if e != nil {
						err = errors.New("parse err")
						break LOOP1
					}

					list[i] |= 1 << ii
					continue
				}
			}
			if err != nil {
				err = errors.New("parse err")
				return
			}

			clock, e := NewClock(interval, list[0], list[1], list[2], list[3], list[4], list[5])
			if e != nil {
				err = e
				return
			}
			j.clock = clock
			now = clock.Now()
		} else {
			err = errors.New("parse err")
			return
		}
	}

	if interval == time.Minute {
		now = time.Unix(now.Unix()-int64(now.Second()), 0)
	}

	j.nextTime = now

	return
}

func (j *Job) Next(interval time.Duration) (slot int, err error) {
	now := j.nextTime
	if j.everyValue > 0 {
		switch j.everyName {
		case "second":
			now = now.Add(time.Second * time.Duration(j.everyValue))
		case "minute":
			now = now.Add(time.Minute * time.Duration(j.everyValue))
		case "hour":
			now = now.Add(time.Hour * time.Duration(j.everyValue))
		case "day":
			now = now.AddDate(0, 0, j.everyValue)
		case "month":
			now = now.AddDate(0, j.everyValue, 0)
		case "week":
			now = now.AddDate(0, 0, 7*j.everyValue)
		}
	} else {
		now, err = j.clock.NextWithWeek()
		if err != nil {
			return
		}
	}

	if interval == time.Minute {
		now = time.Unix(now.Unix()-int64(now.Second()), 0)
	}

	if j.RunIfDelay {
		if now.Sub(time.Now()) < interval {
			now = time.Now().Add(interval)
		}
	} else {
		if now.Sub(time.Now()) < interval {
			return j.Next(interval)
		}
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
