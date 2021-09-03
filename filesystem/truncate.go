package filesystem

import (
	"io"
	"math"
	"os"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"
)

const (
	// DefaultTruncateUnitSize represents default truncate size.
	DefaultTruncateUnitSize = 1024 * 1024
)

// FileTruncator encapsulates necessary data for truncation
type FileTruncator struct {
	FilePath         string
	FileSize         int64
	writer           io.Writer
	TruncateUnit     int64
	TruncateInterval time.Duration
}

// NewFileTruncator returns a new file truncator
func NewFileTruncator(filePath string, interval time.Duration, unitSize int64, writer io.Writer) *FileTruncator {
	truncator := &FileTruncator{
		FilePath:         filePath,
		TruncateUnit:     unitSize,
		TruncateInterval: interval,
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
	truncateCount := truncator.TruncateCount()
	bar := pb.New(truncateCount).SetUnits(pb.MiB)
	bar.Output = truncator.writer
	bar.Start()
	for i := 0; i < truncateCount; i++ {
		bar.Increment()
		time.Sleep(truncator.TruncateInterval)
		file.Truncate(truncator.TruncateUnit)
		file.Sync()
	}
	bar.FinishPrint("Removed " + truncator.FilePath)
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
