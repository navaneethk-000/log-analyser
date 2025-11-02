package filter

import (
	"log_parser/model"
	"testing"
	"time"
)

func parseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}

func sampleLogStore() model.LogStore {
	entries := []model.LogEntry{
		{Log: "entry1", Level: "INFO", Component: "auth", Host: "server1", Request_id: "r1", Time: parseTime("2025-01-01 10:00:00")},
		{Log: "entry2", Level: "ERROR", Component: "db", Host: "server2", Request_id: "r2", Time: parseTime("2025-01-01 11:00:00")},
		{Log: "entry3", Level: "INFO", Component: "api", Host: "server1", Request_id: "r3", Time: parseTime("2025-01-01 12:00:00")},
	}

	index := model.SegmentIndex{
		ByLevel: map[string][]int{
			"INFO":  {0, 2},
			"ERROR": {1},
		},
		ByComponent: map[string][]int{
			"auth": {0},
			"db":   {1},
			"api":  {2},
		},
		ByHost: map[string][]int{
			"server1": {0, 2},
			"server2": {1},
		},
		ByReqID: map[string][]int{
			"r1": {0},
			"r2": {1},
			"r3": {2},
		},
	}

	segment := model.Segment{
		LogEntries: entries,
		StartTime:  parseTime("2025-01-01 09:00:00"),
		EndTime:    parseTime("2025-01-01 13:00:00"),
		Index:      index,
	}

	return model.LogStore{Segments: []model.Segment{segment}}
}

func TestNoFilters(t *testing.T) {
	store := sampleLogStore()

	got := FilterLogs(store, nil, nil, nil, nil, time.Time{}, time.Time{})
	if len(got) != 3 {
		t.Errorf("expected all 3 logs with no filters,but got %d", len(got))
	}
}

func TestFilterByLevel(t *testing.T) {
	store := sampleLogStore()

	got := FilterLogs(store, []string{"ERROR"}, nil, nil, nil, time.Time{}, time.Time{})
	if len(got) != 1 || got[0].Log != "entry2" {
		t.Errorf("expected entry2,but got %+v", got)
	}
}

func TestFilterByComponent(t *testing.T) {
	store := sampleLogStore()

	got := FilterLogs(store, nil, []string{"auth"}, nil, nil, time.Time{}, time.Time{})
	if len(got) != 1 || got[0].Log != "entry1" {
		t.Errorf("expected entry1,but got %+v", got)
	}
}

func TestFilterByHost(t *testing.T) {
	store := sampleLogStore()

	got := FilterLogs(store, nil, nil, []string{"server2"}, nil, time.Time{}, time.Time{})
	if len(got) != 1 || got[0].Log != "entry2" {
		t.Errorf("expected entry2,but got %+v", got)
	}
}

func TestFilterByReqID(t *testing.T) {
	store := sampleLogStore()

	got := FilterLogs(store, nil, nil, nil, []string{"r3"}, time.Time{}, time.Time{})
	if len(got) != 1 || got[0].Log != "entry3" {
		t.Errorf("expected entry3,but got %+v", got)
	}
}

func TestFilterByTimeRange(t *testing.T) {
	store := sampleLogStore()

	start := parseTime("2025-01-01 10:30:00")
	end := parseTime("2025-01-01 12:30:00")
	got := FilterLogs(store, nil, nil, nil, nil, start, end)

	if len(got) != 2 {
		t.Errorf("expected 2 logs between 10:30 and 12:30,but got %d", len(got))
	}
}

func TestBothLevelAndComponent(t *testing.T) {
	store := sampleLogStore()

	got := FilterLogs(store, []string{"INFO"}, []string{"api"}, nil, nil, time.Time{}, time.Time{})
	if len(got) != 1 || got[0].Log != "entry3" {
		t.Errorf("expected entry3, got %+v", got)
	}
}

func TestStoreEmpty(t *testing.T) {
	store := model.LogStore{}
	got := FilterLogs(store, nil, nil, nil, nil, time.Time{}, time.Time{})
	if len(got) != 0 {
		t.Errorf("expected 0 logs for empty store,but got %d", len(got))
	}
}
func TestSkipSegmentBeforeStartTime(t *testing.T) {
	now := time.Now()
	store := model.LogStore{
		Segments: []model.Segment{
			{
				StartTime:  now.Add(-2 * time.Hour),
				EndTime:    now.Add(-1 * time.Hour),
				LogEntries: []model.LogEntry{},
			},
		},
	}
	startTime := now
	endTime := time.Time{}

	result := FilterLogs(store, nil, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs, got %d", len(result))
	}
}
func TestSkipSegmentAfterEndTime(t *testing.T) {
	now := time.Now()
	store := model.LogStore{
		Segments: []model.Segment{
			{
				StartTime:  now.Add(2 * time.Hour),
				EndTime:    now.Add(3 * time.Hour),
				LogEntries: []model.LogEntry{},
			},
		},
	}
	endTime := now
	startTime := time.Time{}

	result := FilterLogs(store, nil, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs, got %d", len(result))
	}
}
func TestSkipEntryBeforeStartTime(t *testing.T) {
	now := time.Now()
	entry := model.LogEntry{Time: now.Add(-2 * time.Hour)}
	segment := model.Segment{
		StartTime:  now.Add(-3 * time.Hour),
		EndTime:    now.Add(-1 * time.Hour),
		LogEntries: []model.LogEntry{entry},
	}
	store := model.LogStore{Segments: []model.Segment{segment}}

	startTime := now
	endTime := time.Time{}

	result := FilterLogs(store, nil, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs, got %d", len(result))
	}
}
func TestSkipEntryAfterEndTime(t *testing.T) {
	now := time.Now()
	entry := model.LogEntry{Time: now.Add(2 * time.Hour)}
	segment := model.Segment{
		StartTime:  now,
		EndTime:    now.Add(3 * time.Hour),
		LogEntries: []model.LogEntry{entry},
	}
	store := model.LogStore{Segments: []model.Segment{segment}}

	startTime := time.Time{}
	endTime := now

	result := FilterLogs(store, nil, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs, got %d", len(result))
	}
}
func TestMatchedEntryBeforeStartTime(t *testing.T) {
	now := time.Now()

	entry1 := model.LogEntry{
		Time: now.Add(-2 * time.Hour),
		Log:  "entry-before-start",
	}
	segment := model.Segment{
		StartTime:  now.Add(-3 * time.Hour),
		EndTime:    now,
		LogEntries: []model.LogEntry{entry1},
		Index: model.SegmentIndex{
			ByLevel: map[string][]int{"INFO": {0}},
		},
	}
	store := model.LogStore{Segments: []model.Segment{segment}}

	startTime := now.Add(-1 * time.Hour)
	endTime := time.Time{}

	result := FilterLogs(store, []string{"INFO"}, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs,but got %d", len(result))
	}
}
func TestMatchedEntryAfterEndTime(t *testing.T) {
	now := time.Now()

	entry1 := model.LogEntry{
		Time: now.Add(2 * time.Hour),
		Log:  "entry-after-end",
	}
	segment := model.Segment{
		StartTime:  now,
		EndTime:    now.Add(3 * time.Hour),
		LogEntries: []model.LogEntry{entry1},
		Index: model.SegmentIndex{
			ByLevel: map[string][]int{"INFO": {0}},
		},
	}
	store := model.LogStore{Segments: []model.Segment{segment}}

	startTime := time.Time{}
	endTime := now.Add(1 * time.Hour)

	result := FilterLogs(store, []string{"INFO"}, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs,but got %d", len(result))
	}
}
