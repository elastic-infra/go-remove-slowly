package filesystem

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/elastic-infra/go-remove-slowly/output"
)

func TestNewFileTruncator(t *testing.T) {
	truncator := NewFileTruncator("path", time.Duration(10), DefaultTruncateSizeMB, output.Type_ProgressBar, nil)
	if truncator.FilePath != "path" {
		t.Fatalf("FilePath is incorrect")
	}
}

func TestTruncateCount(t *testing.T) {
	truncator := NewFileTruncator("path", time.Duration(10), DefaultTruncateSizeMB, output.Type_ProgressBar, nil)
	tests := []struct {
		size  int64
		count int
	}{
		{5 * truncateSizeUnit, 5},
		{5*truncateSizeUnit + 5, 6},
		{5*truncateSizeUnit - 5, 5},
		{512, 1},
	}
	for _, test := range tests {
		answer := test.count
		truncator.FileSize = test.size
		actual := truncator.TruncateCount()
		if actual != answer {
			t.Fatalf("Truncate count is wrong: expected(%d) actual(%d)", answer, actual)
		}
	}
}

func TestUpdateStat(t *testing.T) {
	path := fmt.Sprintf("%s/%s", os.TempDir(), "statTestFile")
	file, err := os.Create(path)
	defer func() { os.Remove(path) }()
	if err != nil {
		t.Fatalf("Failed to create file %s", err.Error())
	}
	var size int64
	size = 10 * truncateSizeUnit
	_, err = file.WriteAt([]byte("a"), size-1)
	if err != nil {
		t.Fatalf("Failed to write to the file %s", err.Error())
	}
	truncator := NewFileTruncator(path, time.Duration(1), DefaultTruncateSizeMB, output.Type_ProgressBar, nil)
	err = truncator.UpdateStat()
	if err != nil {
		t.Fatalf("UpdateStat failed: %s", err.Error())
	}
	if truncator.FileSize != size {
		t.Fatalf("UpdateStat did not correctly get file size expected(%d) actual(%d)", size, truncator.FileSize)
	}
}

func TestUpdateStat_PathError(t *testing.T) {
	path := fmt.Sprintf("%s/%s", os.TempDir(), "statTestFile")
	truncator := NewFileTruncator(path, time.Duration(1), DefaultTruncateSizeMB, output.Type_ProgressBar, nil)
	err := truncator.UpdateStat()
	if err == nil {
		t.Fatal("Error did not happen")
	}
}

func TestRemove(t *testing.T) {
	path := fmt.Sprintf("%s/%s", os.TempDir(), "testRemoveFile")
	file, err := os.Create(path)
	defer func() { os.Remove(path) }()
	if err != nil {
		t.Fatalf("Failed to create file %s", err.Error())
	}
	var size int64
	size = 10 * truncateSizeUnit
	_, err = file.WriteAt([]byte("a"), size-1)
	if err != nil {
		t.Fatalf("Failed to write to the file %s", err.Error())
	}
	truncator := NewFileTruncator(path, time.Duration(1), DefaultTruncateSizeMB, output.Type_ProgressBar, nil)
	err = truncator.Remove()
	if err != nil {
		t.Fatalf("File Removal failed %s", err.Error())
	}
	_, err = os.Stat(path)
	if err == nil {
		t.Fatalf("File %s is not removed", path)
	}
}
