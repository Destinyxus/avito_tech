package utils

import "time"

func SetPeriod(month, year int) (first, second time.Time) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return startDate, endDate
}
