package main

import "io"

// IoMayDumbWriter provides io.Writer interface that does not write when isDumb is true
type IoMayDumbWriter struct {
	stream io.Writer
	isDumb bool
}

// NewIoMayDumbWriter returns a new io.Writer interface of that
func NewIoMayDumbWriter(stream io.Writer, isDumb bool) io.Writer {
	writer := &IoMayDumbWriter{stream, isDumb}
	return writer
}

func (writer *IoMayDumbWriter) Write(s []byte) (int, error) {
	if writer.isDumb {
		return len(s), nil
	}
	return writer.stream.Write(s)
}
