package filter

import (
	"fmt"
	"log_parser/model"
	"sync"
	"time"
)

func FilterLogs(
	store model.LogStore,
	levels, components, hosts, reqIDs []string,
	startTime time.Time, endTime time.Time,
) []model.LogEntry {
	start := time.Now()

	var result []model.LogEntry
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, segment := range store.Segments {
		wg.Add(1)
		go func(seg model.Segment) {
			defer wg.Done()

			totalFilters := 0
			if !startTime.IsZero() && seg.EndTime.Before(startTime) {
				return
			}
			if !endTime.IsZero() && seg.StartTime.After(endTime) {
				return
			}

			matchedIndex := make(map[int]bool)

			if len(levels) > 0 {
				totalFilters++
				for _, level := range levels {
					for _, idx := range seg.Index.ByLevel[level] {
						matchedIndex[idx] = true
					}
				}
			}

			if len(components) > 0 {
				totalFilters++
				componentFilter := make(map[int]bool)
				for _, component := range components {
					for _, idx := range seg.Index.ByComponent[component] {
						if len(matchedIndex) == 0 || matchedIndex[idx] {
							componentFilter[idx] = true
						}
					}
				}
				matchedIndex = componentFilter
			}

			if len(hosts) > 0 {
				totalFilters++
				hostFilter := make(map[int]bool)
				for _, host := range hosts {
					for _, idx := range seg.Index.ByHost[host] {
						if len(matchedIndex) == 0 || matchedIndex[idx] {
							hostFilter[idx] = true
						}
					}
				}
				matchedIndex = hostFilter
			}

			if len(reqIDs) > 0 {
				totalFilters++
				requestFilter := make(map[int]bool)
				for _, reqID := range reqIDs {
					for _, idx := range seg.Index.ByReqID[reqID] {
						if len(matchedIndex) == 0 || matchedIndex[idx] {
							requestFilter[idx] = true
						}
					}
				}
				matchedIndex = requestFilter
			}

			var filteredEntries []model.LogEntry

			if totalFilters == 0 && startTime.IsZero() && endTime.IsZero() {
				filteredEntries = append(filteredEntries, seg.LogEntries...)
			} else if totalFilters == 0 {
				for _, entry := range seg.LogEntries {
					if !startTime.IsZero() && entry.Time.Before(startTime) {
						continue
					}
					if !endTime.IsZero() && entry.Time.After(endTime) {
						continue
					}
					filteredEntries = append(filteredEntries, entry)
				}
			}

			for idx := range matchedIndex {
				entry := seg.LogEntries[idx]
				if !startTime.IsZero() && entry.Time.Before(startTime) {
					continue
				}
				if !endTime.IsZero() && entry.Time.After(endTime) {
					continue
				}
				filteredEntries = append(filteredEntries, entry)
			}
			mu.Lock()
			result = append(result, filteredEntries...)
			mu.Unlock()
		}(segment)
	}

	wg.Wait()

	timeTaken := time.Since(start)
	fmt.Println("Time consumed for filtering :", timeTaken)
	return result
}
