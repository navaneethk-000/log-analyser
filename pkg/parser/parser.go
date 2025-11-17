package parser

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"log_parser/model"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func ParseLogEntry(s string) (*model.LogEntry, error) {
	pattern := `^(?P<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)\s+\|\s+(?P<level>[A-Z]+)\s+\|\s+(?P<component>[\w-]+)\s+\|\s+host=(?P<host>[\w-]+)\s+\|\s+request_id=(?P<request_id>[\w-]+)\s+\|\s+msg="(?P<msg>.*)"$`

	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("Error:%v", err)
	}
	matches := r.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("invalid log format")
	}

	time, err := time.Parse("2006-01-02 15:04:05.000", matches[r.SubexpIndex("time")])
	if err != nil {
		return nil, fmt.Errorf("Error:%v", err)
	}

	return &model.LogEntry{
		Log:        matches[0],
		Time:       time,
		Level:      model.LogLevel(matches[r.SubexpIndex("level")]),
		Component:  matches[r.SubexpIndex("component")],
		Host:       matches[r.SubexpIndex("host")],
		Request_id: matches[r.SubexpIndex("request_id")],
		Message:    matches[r.SubexpIndex("msg")],
	}, nil
}

func ParseLogFiles(folderPath string) ([]model.LogEntry, error) {
	var allEntries []model.LogEntry
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Println("Error Parsing the file", folderPath)
		return nil, fmt.Errorf("failed to read directory : %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		slog.Info("Parsing the file", "Name", file.Name())
		path := filepath.Join(folderPath, file.Name())
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("Skipping file %s due to error: %v\n", path, err)
			continue
		}
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)

		for scanner.Scan() {
			line := scanner.Text()
			entry, err := ParseLogEntry(line)
			if err == nil {
				allEntries = append(allEntries, *entry)
			}
		}

		f.Close()

	}
	return allEntries, nil
}
