// Tracking daily focus with the Forest app.
//
// Organize focus sessions by day to see how close you are to your daily focus
// hour goals.
package main

import (
	"flag"
	"fmt"
)

var hoursPerDay int

func main() {
	// Process flags
	flag.IntVar(&hoursPerDay, "hours", 8, "target number of hours to focus per day: between 0 and 23")
	filePath := flag.String("file", "", "path to .csv file (required)")
	port := flag.String("port", "", "serve http on localhost instead of writing to .csv")
	outPath := flag.String("out", "tracked_focus.csv", "path to output .csv file")
	trackStart := flag.String("begin", "", "beginning of tracked period in the format YYYY-MM-DD")
	trackFinish := flag.String("end", "", "end of tracked period in the format YYYY-MM-DD")
	flag.Parse()
	if hoursPerDay < 0 || hoursPerDay > 23 {
		fmt.Println("-hours should be between 0 and 23")
		return
	}
	if *filePath == "" {
		fmt.Println("-file missing")
		return
	}

	// Read session-level data
	sessions, err := readSessions(*filePath)
	if err != nil {
		fmt.Printf("Failed to read the input file %s: %v\n", *filePath, err)
		return
	}
	if len(sessions) == 0 {
		fmt.Println("input file can't be empty")
		return
	}

	// Calculate FocusedTime for every session
	for _, session := range sessions {
		session.FocusedTime = session.EndTime.Sub(session.StartTime.Time)
	}

	// Track all sessions from the input file if not specified otherwise by the user
	updateTrackPeriod(sessions, *trackStart, *trackFinish)

	// Group session-level data by days
	daysMap, err := groupByDays(sessions)
	if err != nil {
		fmt.Printf("Failed to get day-level data from provided input: %s\n", err)
		return
	}

	// Sort FocusedTime data by day and calculate SumDelta of FocusedTime for every day
	days := prepareDays(daysMap)

	if *port != "" {
		// Visualize the results
		if err := serveData(days, *port); err != nil {
			fmt.Printf("Failed to start server on port %s: %s\n", *port, err)
			return
		}
	} else {
		// Store the results
		if err = writeDays(days, *outPath); err != nil {
			fmt.Printf("Failed to write into the output file %s: %s\n", *outPath, err)
			return
		}
		fmt.Printf("Done writing data to %s.\n", *outPath)
	}
}
