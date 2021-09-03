package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/elastic-infra/go-remove-slowly/filesystem"
	"github.com/urfave/cli/v2"
)

type MyApp struct {
	*cli.App
	stream io.Writer
}

func NewMyApp() *MyApp {
	app := cli.NewApp()
	return &MyApp{app, nil}
}

func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func getFilePaths(paths []string) ([]string, error) {
	var files []string
	checked := map[string]bool{}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		isDir, err := isDirectory(path)
		if err != nil {
			return err
		}
		if isDir {
			return nil
		}
		if _, ok := checked[path]; ok {
			return nil
		}
		files = append(files, path)
		checked[path] = true
		return nil
	}

	for _, path := range paths {
		isDir, err := isDirectory(path)
		if err != nil {
			return nil, err
		}

		if !isDir {
			if _, ok := checked[path]; !ok {
				files = append(files, path)
				checked[path] = true
			}
			continue
		}

		if err := filepath.Walk(path, walkFunc); err != nil {
			return nil, err
		}
	}
	return files, nil
}

// NewApp returns a cli app
func NewApp() *MyApp {
	app := NewMyApp()
	app.Name = "remove-slowly"
	app.Usage = "Remove files slowly"

	app.Action = func(context *cli.Context) error {
		if context.Bool("version") {
			fmt.Printf("%s\n", versionString())
			return nil
		}
		isDumb := context.Bool("quiet")
		stream := NewIoMayDumbWriter(os.Stdout, isDumb)
		app.stream = stream
		errs := []error{}

		targetFilePaths, err := getFilePaths(context.Args().Slice())
		if err != nil {
			return err
		}

		for _, target := range targetFilePaths {
			fmt.Fprintln(app.stream, "Removing File: "+target)
			truncator := filesystem.NewFileTruncator(target, context.Duration("interval"), context.Int64("size"), app.stream)
			if err := truncator.Remove(); err != nil {
				fmt.Fprintf(os.Stderr, "File %s removal error: %s\n", target, err.Error())
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			msg := ""
			for _, err := range errs {
				msg = msg + err.Error() + "\n"
			}
			return errors.New(msg)
		}
		return nil
	}

	app.Flags = []cli.Flag{
		&cli.DurationFlag{
			Name:    "interval",
			Aliases: []string{"i"},
			Usage:   "Interval between truncations",
			Value:   time.Duration(10) * time.Millisecond,
		},
		&cli.BoolFlag{
			Name:    "quiet",
			Aliases: []string{"q"},
			Usage:   "When true, no output is written",
		},
		&cli.BoolFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Show version and build information",
		},
		&cli.Int64Flag{
			Name:    "size",
			Aliases: []string{"s"},
			Usage:   "Truncation size at once[MB]",
			Value:   filesystem.DefaultTruncateSizeMB,
		},
	}
	return app
}
