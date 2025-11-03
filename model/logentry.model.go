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
	Log        string    `json:"log"`
	Time       time.Time `json:"time"`
	Level      LogLevel  `json:"level"`
	Component  string    `json:"component"`
	Host       string    `json:"host"`
	Request_id string    `json:"requiest_id"`
	Message    string    `json:"message"`
}
