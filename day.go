package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gocarina/gocsv"
)

var (
	// Data for processing
	trackStart  time.Time
	trackFinish time.Time
)

// A Day has the total focused time for that day.
type Day struct {
	Day         CustomTime    `csv:"Day"`
	FocusedTime time.Duration `csv:"Focused Time"`
	Delta       time.Duration `csv:"Delta"`
	SumDelta    time.Duration `csv:"Sum Delta"`
}

func isWorkday(day CustomTime) bool {
	if (day.After(trackStart) || day.Equal(trackStart)) &&
		day.Before(trackFinish) {
		return true
	}
	return false
}

func groupByDays(sessions []*Session) (map[CustomTime]*Day, error) {
	// Add every day
	capacity := int(trackFinish.Sub(trackStart).Hours() / 24)
	daysMap := make(map[CustomTime]*Day, capacity)
	for t := trackStart; t.Before(trackFinish); t = t.Add(24 * time.Hour) {
		ct := CustomTime{t}
		daysMap[ct] = &Day{Day: ct}
	}

	// Calculate FocusedTime for every day
	for _, session := range sessions {
		day := session.StartTime.TruncateDay()
		focusedTime := session.FocusedTime

		if isWorkday(day) {
			if _, ok := daysMap[day]; ok {
				daysMap[day].FocusedTime += focusedTime
			} else {
				return nil, fmt.Errorf("failed to process the input file for day %v", day)
			}
		}
	}

	return daysMap, nil
}

func prepareDays(daysMap map[CustomTime]*Day) []*Day {
	// Sort Focused Time data by day
	days := make([]*Day, 0, len(daysMap))
	for _, day := range daysMap {
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool { return days[i].Day.Before(days[j].Day.Time) })

	// Calculate SumDelta of FocusedTime for every day
	var prevDay *Day
	for _, day := range days {
		day.Delta = time.Duration(hoursPerDay)*time.Hour - day.FocusedTime
		if prevDay != nil {
			day.SumDelta = day.Delta + prevDay.SumDelta
		} else {
			day.SumDelta = day.Delta
		}
		prevDay = day
	}

	return days
}

func writeDays(days []*Day, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	return gocsv.Marshal(days, out)
}
