package model

import "time"

type Segment struct {
	FileName   string
	LogEntries []LogEntry
	StartTime  time.Time
	EndTime    time.Time
	Index      SegmentIndex
}

type SegmentIndex struct {
	ByLevel     map[string][]int
	ByComponent map[string][]int
	ByHost      map[string][]int
	ByReqID     map[string][]int
}

type LogStore struct {
	Segments []Segment
}
