package filesystem

import (
	"fmt"
	"io"
	"math"
	"os"
	"time"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/elastic-infra/go-remove-slowly/output"
)

const (
	// DefaultTruncateSizeMB represents default truncate size.
	DefaultTruncateSizeMB = 1
	truncateSizeUnit      = 1024 * 1024 // MB
)

// FileTruncator encapsulates necessary data for truncation
type FileTruncator struct {
	FilePath         string
	FileSize         int64
	writer           io.Writer
	TruncateUnit     int64
	TruncateInterval time.Duration
	OutputType       output.Type
}

// NewFileTruncator returns a new file truncator
func NewFileTruncator(filePath string, interval time.Duration, sizeMB int64, outputType output.Type, writer io.Writer) *FileTruncator {
	truncator := &FileTruncator{
		FilePath:         filePath,
		TruncateUnit:     sizeMB * truncateSizeUnit,
		TruncateInterval: interval,
		OutputType:       outputType,
		writer:           writer,
	}
	return truncator
}

// Remove removes a file after gradually truncating each after configured interval
func (truncator *FileTruncator) Remove() error {
	if err := truncator.UpdateStat(); err != nil {
		return err
	}
	file, err := os.OpenFile(truncator.FilePath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	startTime := time.Now()
	truncateCount := truncator.TruncateCount()

	var bar *pb.ProgressBar
	if truncator.OutputType == output.Type_ProgressBar {
		tmpl := "\r" + `{{with string . "prefix"}}{{.}} {{end}}{{counters . }} {{bar . }} {{percent . }} {{rtime . "ETA %s"}}{{with string . "suffix"}} {{.}}{{end}}`
		bar = pb.ProgressBarTemplate(tmpl).New(truncateCount)
		bar.SetRefreshRate(time.Second)
		bar.SetWriter(truncator.writer)
		bar.Start()
	}

	for i := 0; i < truncateCount; i++ {
		if truncator.OutputType == output.Type_ProgressBar {
			bar.Increment()
		}

		var eta time.Duration
		if i != 0 {
			eta = time.Duration(int(time.Since(startTime)) * (truncateCount - i) / i)
		}

		if truncator.OutputType == output.Type_Simple {
			fmt.Printf("file: %s\tcompletion: %0.2f%%\tremaining: %s\n",
				truncator.FilePath,
				float64(i)/float64(truncateCount)*100,
				eta.Round(time.Millisecond),
			)
		}

		time.Sleep(truncator.TruncateInterval)
		file.Truncate(truncator.FileSize - int64(i)*truncator.TruncateUnit)
		file.Sync()
	}

	if truncator.OutputType == output.Type_ProgressBar {
		bar.Finish()
		finishMessage := fmt.Sprintf("Removed " + truncator.FilePath)
		fmt.Println(finishMessage)
	}

	file.Close()
	if err := os.Remove(truncator.FilePath); err != nil {
		return err
	}
	return nil
}

// UpdateStat updates stat information such as FileSize
func (truncator *FileTruncator) UpdateStat() error {
	fileInfo, err := os.Stat(truncator.FilePath)
	if err != nil {
		return err
	}
	truncator.FileSize = fileInfo.Size()
	return nil
}

// TruncateCount returns how many times Truncate() will be called
func (truncator *FileTruncator) TruncateCount() int {
	return int(math.Ceil(float64(truncator.FileSize) / float64(truncator.TruncateUnit)))
}
