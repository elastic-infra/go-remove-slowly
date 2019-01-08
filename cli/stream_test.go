package main

import (
	"strings"
	"testing"
)

func TestNewIoMayDumbWriter(t *testing.T) {
	message := "Hello!"
	bytesMessage := []byte(message)
	tests := []struct {
		isDumb bool
		answer string
	}{{true, ""}, {false, message}}
	for _, test := range tests {
		bufWriter := &strings.Builder{}
		writer := NewIoMayDumbWriter(bufWriter, test.isDumb)
		written, err := writer.Write(bytesMessage)
		if err != nil {
			t.Fatalf("Write() error: %s", err.Error())
		}
		if written != len(message) {
			t.Fatalf("Write() returned wrong written bytes: expected %d actual %d", len(message), written)
		}
		writtenString := bufWriter.String()
		if writtenString != test.answer {
			t.Fatalf("Write did not write the given message: expected(%s), actual(%s)", test.answer, writtenString)
		}
	}
}
