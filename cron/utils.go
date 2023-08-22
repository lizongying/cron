package cron

import (
	"math"
	"time"
)

func LastMoment(interval time.Duration) (lastMoment uint32) {
	if time.Now().Year()%4 == 0 {
		lastMoment = 366 * 24 * 60
	} else {
		lastMoment = 365 * 24 * 60
	}
	if interval == time.Second {
		lastMoment *= 60
	}
	lastMoment--
	return
}

func SlotSinceYear(now time.Time, interval time.Duration) (slot uint32) {
	year, _ := time.ParseInLocation("2006", now.Format("2006"), time.Local)
	if interval == time.Minute {
		slot = uint32(math.Floor(now.Sub(year).Minutes()))
		return
	}

	slot = uint32(math.Floor(now.Sub(year).Seconds()))
	return
}
