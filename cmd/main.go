package main

import (
	"flag"
	"fmt"
	"log/slog"
	"log_parser/filter"
	"log_parser/model"
	"log_parser/segmenter"
	"strings"
	"time"
)

func main() {

	level := flag.String("level", "", "Filter by log level")
	component := flag.String("component", "", "Filter by component")
	host := flag.String("host", "", "Filter by host")
	reqID := flag.String("reqID", "", "Filter by requestID")
	startTimeString := flag.String("after", "", "Filter by start time")
	endTimeString := flag.String("before", "", "Filter by end time")
	flag.Parse()

	logStore, err := segmenter.GenerateSegments("/home/navaneeth/project_log_parser/logs")
	if err != nil {
		slog.Error("Failed to parse logs\n")
	}
	split := func(s string) []string {
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	levels := split(*level)
	components := split(*component)
	hosts := split(*host)
	reqIDs := split(*reqID)
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
