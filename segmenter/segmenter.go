package segmenter

import (
	"bufio"
	"fmt"
	"log_parser/indexer"
	"log_parser/model"
	"log_parser/parser"
	"os"
	"path/filepath"
)

func GenerateSegments(folderPath string) ([]model.Segment, error) {
	var segments []model.Segment

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(folderPath, file.Name())
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()

		var entries []model.LogEntry
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			entry, err := parser.ParseLogEntry(line)
			if err == nil && entry != nil {
				entries = append(entries, *entry)
			}
		}

		if len(entries) == 0 {
			continue
		}

		segment := model.Segment{
			FileName:   file.Name(),
			LogEntries: entries,
			StartTime:  entries[0].Time,
			EndTime:    entries[len(entries)-1].Time,
			Index:      indexer.CreateSegmentIndex(entries),
		}

		segments = append(segments, segment)
	}

	return segments, nil
}
