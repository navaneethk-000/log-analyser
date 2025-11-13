package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"log_parser/segmenter"
	"os"
)

func main() {
	LogfolderPath := flag.String("path", "/home/navaneeth/project_log_parser/logs", "path of log folder")
	outputFile := flag.String("out", "parsed_logs.json", "Output JSON file path")

	flag.Parse()

	segments, err := segmenter.GenerateSegments(*LogfolderPath)
	if err != nil {
		slog.Error("Failed to parse logs", "error", err)
	}

	file, err := os.Create(*outputFile)
	if err != nil {
		slog.Error("Could not create output file", "error", err)

	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(segments, "", " ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("Logs parsed successfully and saved to %s\n", *outputFile)
}
	