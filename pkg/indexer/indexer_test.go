package indexer

import (
	"log_parser/model"
	"reflect"
	"testing"
	"time"
)

func TestCreateSegmentIndex(t *testing.T) {
	entries := []model.LogEntry{
		{
			Time:       time.Now(),
			Level:      model.INFO,
			Component:  "api-server",
			Host:       "worker01",
			Request_id: "req-1",
			Message:    "started service",
		},
		{
			Time:       time.Now(),
			Level:      model.ERROR,
			Component:  "db",
			Host:       "worker02",
			Request_id: "req-2",
			Message:    "database connection failed",
		},
		{
			Time:       time.Now(),
			Level:      model.INFO,
			Component:  "api-server",
			Host:       "worker01",
			Request_id: "req-3",
			Message:    "processed request",
		},
	}

	index := CreateSegmentIndex(entries)

	// ByLevel
	if got, want := index.ByLevel["INFO"], []int{0, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByLevel[INFO] = %v, want %v", got, want)
	}
	if got, want := index.ByLevel["ERROR"], []int{1}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByLevel[ERROR] = %v, want %v", got, want)
	}

	//ByComponent
	if got, want := index.ByComponent["api-server"], []int{0, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByComponent[api-server] = %v, want %v", got, want)
	}
	if got, want := index.ByComponent["db"], []int{1}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByComponent[db] = %v, want %v", got, want)
	}

	//ByHost
	if got, want := index.ByHost["worker01"], []int{0, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByHost[worker01] = %v, want %v", got, want)
	}
	if got, want := index.ByHost["worker02"], []int{1}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByHost[worker02] = %v, want %v", got, want)
	}

	//ByReqID
	if got, want := index.ByReqID["req-1"], []int{0}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByReqID[req-1] = %v, want %v", got, want)
	}
	if got, want := index.ByReqID["req-3"], []int{2}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByReqID[req-3] = %v, want %v", got, want)
	}
}
