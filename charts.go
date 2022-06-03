package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

const (
	chartTimestampLayout = "2006-01-02"
)

var (
	// Data for visualization
	xAxis           []string
	focusedTimeData []opts.LineData
	deltaData       []opts.LineData
	sumDeltaData    []opts.LineData
	targetTimeData  []opts.LineData
	dailyChart      *charts.Line
	sumChart        *charts.Line
)

func serveData(days []*Day, port string) error {
	prepareCharts(days)

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveSumChart)
	mux.HandleFunc("/daily", serveDailyChart)
	fmt.Printf("Starting server on port %s...\n", port)
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return server.ListenAndServe()
}

func prepareCharts(days []*Day) {
	var (
		deltaCategory = fmt.Sprintf("Focused Time Delta (%dh minus Focused Time)", hoursPerDay)
		theme         = types.ThemeInfographic
	)

	for _, day := range days {
		xAxis = append(xAxis, day.Day.Format(chartTimestampLayout))
		focusedTimeData = append(focusedTimeData, opts.LineData{Value: day.FocusedTime.Hours()})
		deltaData = append(deltaData, opts.LineData{Value: day.Delta.Hours()})
		sumDeltaData = append(sumDeltaData, opts.LineData{Value: day.SumDelta.Hours()})
		targetTimeData = append(targetTimeData, opts.LineData{Value: hoursPerDay})
	}

	dailyChart = charts.NewLine()
	dailyChart.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: theme}),
		charts.WithTitleOpts(opts.Title{
			Title:    "How focused you were over time (daily)",
			Subtitle: "How many hours you focused every day.",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:   true,
			Orient: "horizontal",
			Bottom: "0",
			Selected: map[string]bool{
				deltaCategory: false,
			},
		}),
	)
	dailyChart.SetXAxis(xAxis).
		AddSeries("Target Focused Time", targetTimeData).
		AddSeries("Actual Focused Time", focusedTimeData).
		AddSeries(deltaCategory, deltaData)

	sumChart = charts.NewLine()
	sumChart.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: theme}),
		charts.WithTitleOpts(opts.Title{
			Title:    "How focused you were over time (running total)",
			Subtitle: "Above zero: you underperform. Under zero: you overperform.",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:   true,
			Orient: "horizontal",
			Bottom: "0",
		}),
	)
	sumChart.SetXAxis(xAxis).
		AddSeries("How far you are from your goal (in hours)", sumDeltaData)
}

func serveDailyChart(w http.ResponseWriter, _ *http.Request) {
	dailyChart.Render(w)
}

func serveSumChart(w http.ResponseWriter, _ *http.Request) {
	sumChart.Render(w)
}
