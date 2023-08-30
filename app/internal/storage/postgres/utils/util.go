package utils

import (
	"fmt"
	"time"
)

func SetPeriod(month, year int) (first, second time.Time) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return startDate, endDate
}

func SetTTL(day time.Duration) (time.Time, error) {
	if day == 0 {
		return time.Time{}, fmt.Errorf("no ttl")
	}
	current := time.Now()
	ttl := current.Add(time.Hour * 24 * day).Add(-time.Nanosecond)

	return ttl, nil
}
