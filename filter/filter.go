package filter

import (
	"log_parser/model"
	"time"
)

func FilterLogs(
	store model.LogStore,
	levels, components, hosts, reqIDs []string,
	startTime time.Time, endTime time.Time,
) []model.LogEntry {
	var result []model.LogEntry

	for _, segment := range store.Segments {
		totalFilters := 0
		if !startTime.IsZero() && segment.EndTime.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && segment.StartTime.After(endTime) {
			continue
		}
		matchedIndex := make(map[int]bool)
		if len(levels) > 0 {
			totalFilters++
			for _, level := range levels {
				for _, idx := range segment.Index.ByLevel[level] {
					matchedIndex[idx] = true
				}
			}
		}
		if len(components) > 0 {
			totalFilters++
			componentFilter := make(map[int]bool)
			for _, component := range components {
				for _, idx := range segment.Index.ByComponent[component] {
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
				for _, idx := range segment.Index.ByHost[host] {
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
				for _, idx := range segment.Index.ByReqID[reqID] {
					if len(matchedIndex) == 0 || matchedIndex[idx] {
						requestFilter[idx] = true
					}
				}
			}
			matchedIndex = requestFilter
		}
		if totalFilters == 0 && startTime.IsZero() && endTime.IsZero() {
			result = append(result, segment.LogEntries...)
		} else if totalFilters == 0 {
			for _, entry := range segment.LogEntries {
				if !startTime.IsZero() && entry.Time.Before(startTime) {
					continue
				}
				if !endTime.IsZero() && entry.Time.After(endTime) {
					continue
				}
				result = append(result, entry)
			}
		}
		for idx := range matchedIndex {
			entry := segment.LogEntries[idx]
			if !startTime.IsZero() && entry.Time.Before(startTime) {
				continue
			}
			if !endTime.IsZero() && entry.Time.After(endTime) {
				continue
			}
			result = append(result, entry)
		}

	}

	return result
}
