package main

import (
	"time"
)

const (
	inputTimestampLayout = "2006-01-02T15:04:05.000-0700"
)

// A CustomTime is an extension of time.Time with helper methods.
type CustomTime struct {
	time.Time
}

// UnmarshalCSV parses the data from CSV and stores the result in t.
func (t *CustomTime) UnmarshalCSV(csv string) (err error) {
	t.Time, err = time.Parse(inputTimestampLayout, csv)
	t.Time = t.Time.UTC()
	return
}

// TruncateDay returns the result of rounding t down toward the closest day.
func (t CustomTime) TruncateDay() CustomTime {
	return CustomTime{time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())}
}
