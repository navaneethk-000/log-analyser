package helper

import (
	"log/slog"
	"strings"
	"time"
)

func Split(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func ParseTimeFlag(timeString string) time.Time {
	var parsedTime time.Time
	if timeString != "" {
		var err error
		parsedTime, err = time.Parse("2006-01-02 15:04:05", timeString)
		if err != nil {
			slog.Error("Error parsing start time", "error", err)
		}
	}
	return parsedTime
}
