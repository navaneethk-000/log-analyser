package segmenter

import (
	"bufio"
	"fmt"
	"log_parser/model"
	"log_parser/pkg/indexer"
	"log_parser/pkg/parser"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func GenerateSegments(folderPath string) ([]model.Segment, error) {
	start := time.Now()
	var segments []model.Segment
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()

			path := filepath.Join(folderPath, file.Name())
			f, err := os.Open(path)
			if err != nil {
				return
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
				return
			}

			segment := model.Segment{
				FileName:   file.Name(),
				LogEntries: entries,
				StartTime:  entries[0].Time,
				EndTime:    entries[len(entries)-1].Time,
				Index:      indexer.CreateSegmentIndex(entries),
			}

			mu.Lock()
			segments = append(segments, segment)
			mu.Unlock()
		}(file)
	}

	wg.Wait()
	timeTaken := time.Since(start)
	fmt.Println("Time consumed for Segmenting :", timeTaken)

	return segments, nil
}
