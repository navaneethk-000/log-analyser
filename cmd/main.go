package main

import (
	"fmt"
	"log_parser/segmenter"
)

func main() {

	// entries, err := parser.ParseLogFiles("/home/navaneeth/project_log_parser/logs")
	// if err != nil {
	// 	fmt.Printf("Error parsing log files: %v\n", err)
	// 	return
	// }
	// fmt.Printf("\nParsed %d log entries:\n\n", len(entries))
	// for i, entry := range entries {
	// 	fmt.Printf("[%d] Time: %v | Level: %v | Component: %v | Host: %v | ReqID: %v | Msg: %v\n",
	// 		i+1, entry.Time, entry.Level, entry.Component, entry.Host, entry.Request_id, entry.Message)
	// }

	segments, _ := segmenter.GenerateSegments("/home/navaneeth/project_log_parser/logs")
	fmt.Println("Firsts entry : ", segments[0].LogEntries[0].Log)
	fmt.Println("IndexByLevel  : ", segments[0].Index.ByLevel)
}
