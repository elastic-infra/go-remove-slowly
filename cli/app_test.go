package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewMyApp(t *testing.T) {
	app := NewMyApp()
	if app.stream != nil {
		t.Fatalf("stream should be initialized with nil")
	}
}

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app.stream != nil {
		t.Fatalf("stream should be initialized with nil")
	}
}

func TestAction(t *testing.T) {
	app := NewApp()
	tests := []struct {
		path string
		size int64
	}{
		{"fileA", 10 * 1024 * 1024},
		{"fileB", 10*1024 + 100},
		{"fileC", 512},
	}

	tmpfiles := []string{}
	for _, test := range tests {
		path := fmt.Sprintf("%s/%s", os.TempDir(), test.path)
		createFile(path, test.size)
		tmpfiles = append(tmpfiles, path)
		defer func() { os.Remove(path) }()
	}
	err := app.Run(append([]string{"cmd", "-q"}, tmpfiles...))
	if err != nil {
		t.Fatalf("Error happened: %s", err.Error())
	}
	for _, f := range tmpfiles {
		_, perr := os.Stat(f)
		if perr != nil {
			if _, ok := perr.(*os.PathError); !ok {
				t.Fatalf("Filetest for %s failed with non-PathError: %s", f, perr.Error())
			}
		} else {
			t.Fatalf("File %s is not removed", f)
		}
	}
}

func TestAction_PathError(t *testing.T) {
	app := NewApp()
	tests := []struct {
		path string
		size int64
	}{
		{"fileD", 512},
	}

	tmpfiles := []string{}
	for _, test := range tests {
		path := fmt.Sprintf("%s/%s", os.TempDir(), test.path)
		// not creating file
		tmpfiles = append(tmpfiles, path)
		defer func() { os.Remove(path) }()
	}
	err := app.Run(append([]string{"cmd", "-q"}, tmpfiles...))
	if err == nil {
		t.Fatalf("Error did not happen (it should)")
	}
}

func TestAction_Version(t *testing.T) {
	app := NewApp()
	output := captureOutput(func() {
		err := app.Run([]string{"cmd", "-v"})
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
	t.Logf("Output: %s", output)
	if !strings.HasPrefix(output, "go-remove-slowly ") {
		t.Fatalf("version string should begin with the app name")
	}
}

func TestSimpleOutput(t *testing.T) {
	app := NewApp()
	createFile("test/bar/file0", 4*1024*1024)
	output := captureOutput(func() {
		err := app.Run([]string{"cmd", "--output", "simple", "--size", "1", "test/bar/file0"})
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
	t.Logf("Output: %s", output)
	if !strings.Contains(output, "file: test/bar/file0") {
		t.Fatalf("output must contain the name of the file which is being removed")
	}
	if !strings.Contains(output, "completion: 0.00%") {
		t.Fatalf("output must contain 4 lines with each completion percentage ratio")
	}
	if !strings.Contains(output, "remaining: ") {
		t.Fatalf("output must contain the remaining amount of time")
	}
}

func createFile(path string, size int64) {
	file, err := os.Create(path)
	if err != nil {
		panic("Failed to create file " + path + " " + err.Error())
	}
	_, err = file.WriteAt([]byte("a"), size-1)
	if err != nil {
		panic("Failed to write file " + path + " " + err.Error())
	}
}

func Test_getFilePaths(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "file only",
			args: args{
				paths: []string{"test/file0"},
			},
			want: []string{
				"test/file0",
			},
		},
		{
			name: "dup files",
			args: args{
				paths: []string{"test/file0", "test/file0"},
			},
			want: []string{
				"test/file0",
			},
		},
		{
			name: "dir only",
			args: args{
				paths: []string{"test/foo/"},
			},
			want: []string{
				"test/foo/file1",
				"test/foo/file2",
			},
		},
		{
			name: "dup dir",
			args: args{
				paths: []string{"test/", "test/foo/"},
			},
			want: []string{
				"test/file0",
				"test/foo/file1",
				"test/foo/file2",
			},
		},
		{
			name: "file and dir",
			args: args{
				paths: []string{"test/file0", "test/foo/"},
			},
			want: []string{
				"test/file0",
				"test/foo/file1",
				"test/foo/file2",
			},
		},
		{
			name: "no such file or dir",
			args: args{
				paths: []string{"test/xxxx"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFilePaths(tt.args.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFilePaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFilePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isDirectory(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "yes",
			args: args{
				path: "test",
			},
			want: true,
		},
		{
			name: "not",
			args: args{
				path: "main.go",
			},
			want: false,
		},
		{
			name: "error",
			args: args{
				path: "no_path",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isDirectory(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("isDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}
