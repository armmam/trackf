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
	// First assume we track by the whole time period
	trackStart = sessions[0].StartTime.TruncateDay().Time
	trackFinish = sessions[len(sessions)-1].EndTime.TruncateDay().Add(24 * time.Hour)

	// Then apply the flags wherever necessary
	inputStart, errStart := time.Parse(trackPeriodLayout, start)
	inputFinish, errFinish := time.Parse(trackPeriodLayout, finish)

	// If both flags are valid, apply them both
	if errStart == nil && errFinish == nil && inputFinish.Before(inputFinish) {
		trackStart = inputStart
		trackFinish = inputFinish
	}
	// If only one of the flags is valid, apply the valid one
	if errStart == nil && errFinish != nil && inputFinish.Before(trackFinish) {
		trackStart = inputStart
	}
	if errStart != nil && errFinish == nil && trackStart.Before(inputFinish) {
		trackFinish = inputFinish
	}
}
