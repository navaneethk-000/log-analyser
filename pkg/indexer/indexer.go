package indexer

import "log_parser/model"

func CreateSegmentIndex(entries []model.LogEntry) model.SegmentIndex {
	index := model.SegmentIndex{
		ByLevel:     make(map[string][]int),
		ByComponent: make(map[string][]int),
		ByHost:      make(map[string][]int),
		ByReqID:     make(map[string][]int),
	}

	for i, entry := range entries {
		index.ByLevel[string(entry.Level)] = append(index.ByLevel[string(entry.Level)], i)
		index.ByComponent[entry.Component] = append(index.ByComponent[entry.Component], i)
		index.ByHost[entry.Host] = append(index.ByHost[entry.Host], i)
		index.ByReqID[entry.Request_id] = append(index.ByReqID[entry.Request_id], i)
	}
	return index
}
