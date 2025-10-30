package segmenter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func createTestFile(entry []string) (string, error) {
	f, err := os.CreateTemp("", "file1")
	if err != nil {
		fmt.Println("Couldn't create temp file.")
	}
	data := strings.Join(entry, "\n")
	_, err = f.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return filepath.Dir(f.Name()), nil
}

func TestCreateSegment(t *testing.T) {
	entry := []string{
		`2025-10-29 15:40:47.814 | DEBUG | auth | host=db01 | request_id=req-rrfym6-7292 | msg="Session terminated"`,
		`2025-10-29 15:40:48.118 | ERROR | worker | host=web02 | request_id=req-rrfym6-7292 | msg="Transaction rolled back"`, `2025-10-23 17:03:52.202 | INFO | worker | host=web01 | request_id=req-68rw7u-4614 | msg="Transaction rolled back"`,
		`2025-10-29 15:40:48.564 | INFO | cache | host=cache01 | request_id=req-rrfym6-7292 | msg="Job retry scheduled"`,
		`2025-10-29 15:41:28.985 | INFO | api-server | host=web02 | request_id=req-zrt358-693 | msg="User login successful"`,
	}
	path, _ := createTestFile(entry)

	got, err := GenerateSegments(path)
	expectedStartTime, _ := time.Parse("2006-01-02 15:04:05.000", "2025-10-23 16:53:00.033")
	expectedEndTime, _ := time.Parse("2006-01-02 15:04:05.000", "2025-10-23 17:04:14.387")

	if err != nil {
		t.Errorf("Expected Error to be nil, but got: %v", err)
	}

	if !strings.HasPrefix(got[0].FileName, "file1") {
		t.Errorf("Expected filename to start with file1, but got:%v", got[0].FileName)
	}

	if got[0].StartTime != expectedStartTime {
		t.Errorf("Expected start time: %v, but got: %v", expectedStartTime, got[0].StartTime)
	}

	if got[0].EndTime != expectedEndTime {
		t.Errorf("Expected end time: %v, but got: %v", expectedEndTime, got[0].EndTime)
	}

}
