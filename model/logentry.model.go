package model

import "time"

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	DEBUG LogLevel = "DEBUG"
	ERROR LogLevel = "ERROR"
)

type LogEntry struct {
	Log        string
	Time       time.Time
	Level      LogLevel
	Component  string
	Host       string
	Request_id string
	Message    string
}
