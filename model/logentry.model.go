package model

import "time"

type LogLevel string

type LogEntry struct {
	Log        string
	Time       time.Time
	Level      LogLevel
	Component  string
	Host       string
	Request_id string
	Message    string
}
