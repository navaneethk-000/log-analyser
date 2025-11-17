package main

import (
	"flag"
	"fmt"
	"log/slog"
	"log_parser/cmd/logFilter/helper"
	"log_parser/model"
	"log_parser/pkg/filter"
	"log_parser/pkg/segmenter"
	"time"
)

func main() {
	level := flag.String("level", "", "Filter by log level")
	component := flag.String("component", "", "Filter by component")
	host := flag.String("host", "", "Filter by host")
	reqID := flag.String("reqID", "", "Filter by requestID")
	startTimeString := flag.String("after", "", "Filter by start time")
	endTimeString := flag.String("before", "", "Filter by end time")
	folderPath := flag.String("path", "/home/navaneeth/project_log_parser/logs", "path of log folder")
	flag.Parse()

	logStore, err := segmenter.GenerateSegments(*folderPath)

	if err != nil {
		slog.Error("Failed to parse logs\n")
	}

	levels := helper.Split(*level)
	components := helper.Split(*component)
	hosts := helper.Split(*host)
	reqIDs := helper.Split(*reqID)

	var startTime, endTime time.Time

	if *startTimeString != "" {
		startTime, err = time.Parse("2006-01-02 15:04:05", *startTimeString)
		if err != nil {
			slog.Error("Error parsing start time", "error", err)
		}
	}

	if *endTimeString != "" {
		endTime, err = time.Parse("2006-01-02 15:04:05", *endTimeString)
		if err != nil {
			slog.Error("Error parsing end time", "error", err)
		}
	}

	filteredLogs := filter.FilterLogs(model.LogStore{Segments: logStore}, levels, components, hosts, reqIDs, startTime, endTime)
	fmt.Printf("Found %d matching entries\n", len(filteredLogs))

	for _, entry := range filteredLogs {
		fmt.Println(entry.Log)
	}

}
