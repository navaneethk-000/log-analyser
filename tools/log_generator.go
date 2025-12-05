package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	sizeLimit  = flag.Int("size", 1024*1024, "Size limit per log file in bytes")
	numFiles   = flag.Int("files", 5, "Number of log files to generate")
	components = flag.String("components", "api-server,database,auth,cache,worker", "Comma-separated list of components")
	hosts      = flag.String("hosts", "web01,web02,db01,cache01,worker01", "Comma-separated list of hosts")
	outputDir  = flag.String("output", "./logs", "Output directory for log files")
)

var logLevels = []string{"DEBUG", "INFO", "WARN", "ERROR"}

var messageTemplates = map[string][]string{
	"api-server": {
		"Request started GET /api/users",
		"Request started POST /api/orders",
		"Request completed in %dms",
		"Connection pool exhausted",
		"Rate limit exceeded for client",
		"Authentication successful",
		"Authentication failed",
		"Request validation failed",
		"Response sent with status %d",
	},
	"database": {
		"Query executed in %dms",
		"Connection established to replica",
		"Transaction committed",
		"Transaction rolled back",
		"Deadlock detected, retrying",
		"Index scan on users table",
		"Full table scan warning",
		"Connection pool size: %d",
	},
	"auth": {
		"User login successful",
		"User login failed - invalid credentials",
		"Token generated for session",
		"Token validation successful",
		"Token expired",
		"Password reset requested",
		"2FA verification completed",
		"Session terminated",
	},
	"cache": {
		"Cache hit for key %s",
		"Cache miss for key %s",
		"Cache entry expired",
		"Cache cleared",
		"Memory usage at %d%%",
		"Eviction policy triggered",
		"Connection to cache server established",
	},
	"worker": {
		"Job started: process-payment",
		"Job completed: send-email",
		"Job failed: generate-report",
		"Job retry scheduled",
		"Queue length: %d",
		"Worker pool size increased to %d",
		"Background task scheduled",
	},
}

type LogEntry struct {
	timestamp time.Time
	level     string
	component string
	host      string
	requestID string
	message   string
}

func (l LogEntry) String() string {
	return fmt.Sprintf("%s | %s | %s | host=%s | request_id=%s | msg=\"%s\"\n",
		l.timestamp.Format("2006-01-02 15:04:05.000"),
		l.level,
		l.component,
		l.host,
		l.requestID,
		l.message)
}

func generateMessage(component string) string {
	templates, ok := messageTemplates[component]
	if !ok {
		templates = []string{"Generic log message"}
	}

	template := templates[rand.Intn(len(templates))]

	if strings.Contains(template, "%d") {
		return fmt.Sprintf(template, rand.Intn(1000)+1)
	}
	if strings.Contains(template, "%s") {
		keys := []string{"user:123", "session:abc", "product:xyz", "order:456"}
		return fmt.Sprintf(template, keys[rand.Intn(len(keys))])
	}

	return template
}

func generateRequestID() string {
	return fmt.Sprintf("req-%s-%d", randomString(6), rand.Intn(10000))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	componentList := strings.Split(*components, ",")
	hostList := strings.Split(*hosts, ",")

	for i := range componentList {
		componentList[i] = strings.TrimSpace(componentList[i])
	}
	for i := range hostList {
		hostList[i] = strings.TrimSpace(hostList[i])
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generating %d log files with max size %d bytes each...\n", *numFiles, *sizeLimit)
	fmt.Printf("Components: %v\n", componentList)
	fmt.Printf("Hosts: %v\n", hostList)

	// Calculate how many total entries we need based on size limit
	// Add 20% buffer to ensure we have enough entries
	avgEntrySize := 118 // actual average from your output
	totalEntriesNeeded := ((*numFiles * *sizeLimit) * 12) / (avgEntrySize * 10)

	// Generate request timelines
	numRequests := totalEntriesNeeded / 5 // Each request has ~5 entries on average
	if numRequests < 200 {
		numRequests = 200
	}

	requestTimelines := make(map[string][]LogEntry)
	baseTime := time.Now().Add(-24 * time.Hour)

	for i := 0; i < numRequests; i++ {
		requestID := generateRequestID()
		// Each request has 3-8 log entries
		numEntries := rand.Intn(6) + 3
		entries := make([]LogEntry, numEntries)

		reqStartTime := baseTime.Add(time.Duration(rand.Intn(86400)) * time.Second)

		for j := 0; j < numEntries; j++ {
			entries[j] = LogEntry{
				timestamp: reqStartTime.Add(time.Duration(j*rand.Intn(500)+10) * time.Millisecond),
				level:     logLevels[rand.Intn(len(logLevels))],
				component: componentList[rand.Intn(len(componentList))],
				host:      hostList[rand.Intn(len(hostList))],
				requestID: requestID,
				message:   generateMessage(componentList[rand.Intn(len(componentList))]),
			}
		}

		requestTimelines[requestID] = entries
	}

	// Flatten all entries
	allEntries := make([]LogEntry, 0)
	for _, entries := range requestTimelines {
		allEntries = append(allEntries, entries...)
	}

	fmt.Printf("Generated %d log entries to distribute\n", len(allEntries))

	// Shuffle entries to distribute across files
	rand.Shuffle(len(allEntries), func(i, j int) {
		allEntries[i], allEntries[j] = allEntries[j], allEntries[i]
	})

	// Sort all entries by timestamp to maintain temporal ordering
	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].timestamp.Before(allEntries[j].timestamp)
	})

	// Distribute entries across files
	entryIndex := 0
	filesCreated := 0

	for fileNum := 0; fileNum < *numFiles; fileNum++ {
		if entryIndex >= len(allEntries) {
			fmt.Printf("Warning: Ran out of entries after %d files. Generate more requests.\n", filesCreated)
			break
		}

		filename := filepath.Join(*outputDir, fmt.Sprintf("log%04d.log", fileNum+1))

		file, err := os.Create(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", filename, err)
			continue
		}

		var currentSize int64
		entryCount := 0

		// Write entries until we hit size limit
		for entryIndex < len(allEntries) && currentSize < int64(*sizeLimit) {
			entry := allEntries[entryIndex]
			logLine := entry.String()

			n, err := file.WriteString(logLine)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to file %s: %v\n", filename, err)
				break
			}

			currentSize += int64(n)
			entryCount++
			entryIndex++
		}

		file.Close()
		filesCreated++
		fmt.Printf("Generated %s: %d bytes, %d entries\n", filename, currentSize, entryCount)
	}

	fmt.Printf("\nLog generation complete! Created %d files with %d total entries.\n", filesCreated, entryIndex)
}
