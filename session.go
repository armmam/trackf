package main

import (
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

const (
	trackPeriodLayout = "2006-01-02"
)

// A Session is an uninterrupted time period when one focuses.
type Session struct {
	StartTime   CustomTime    `csv:"Start Time"`
	EndTime     CustomTime    `csv:"End Time"`
	FocusedTime time.Duration `csv:"-"`
}

func readSessions(name string) ([]*Session, error) {
	in, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	var sessions []*Session
	return sessions, gocsv.Unmarshal(in, &sessions)
}

func updateTrackPeriod(sessions []*Session, start, finish string) {
	inputStart, errStart := time.Parse(trackPeriodLayout, start)
	inputFinish, errFinish := time.Parse(trackPeriodLayout, finish)

	if errStart != nil && errFinish != nil || !inputStart.Before(inputFinish) {
		// Invalid input, use the whole period
		trackStart = sessions[0].StartTime.TruncateDay().Time
		trackFinish = sessions[len(sessions)-1].EndTime.TruncateDay().Add(24 * time.Hour)
		return
	}

	trackStart = inputStart
	trackFinish = inputFinish
}
