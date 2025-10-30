package model

import "time"

type Segment struct {
	FileName   string
	LogEntries []LogEntry
	StartTime  time.Time
	EndTime    time.Time
}

type LogStore struct {
	Segments []Segment
}
