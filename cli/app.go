package main

import (
	"errors"
	"fmt"
	"io"
	"os"
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

// NewApp returns a cli app
func NewApp() *MyApp {
	app := NewMyApp()
	app.Name = "remove-slowly"
	app.Usage = "Remove files slowly"

	app.Action = func(context *cli.Context) error {
		isDumb := context.Bool("quiet")
		stream := NewIoMayDumbWriter(os.Stdout, isDumb)
		app.stream = stream
		errs := []error{}
		for _, filePath := range context.Args().Slice() {
			fmt.Fprintln(app.stream, "Removing File: "+filePath)
			truncator := filesystem.NewFileTruncator(filePath, context.Duration("interval"), app.stream)
			if err := truncator.Remove(); err != nil {
				fmt.Fprintf(os.Stderr, "File %s removal error: %s\n", filePath, err.Error())
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
			Name:  "interval",
			Aliases: []string{"i"},
			Usage: "Interval between truncations",
			Value: time.Duration(10) * time.Millisecond,
		},
		&cli.BoolFlag{
			Name:  "quiet",
			Aliases: []string{"q"},
			Usage: "When true, no output is written",
		},
		&cli.BoolFlag{
			Name: "version",
			Aliases: []string{"v"},
			Usage: "Show version and build information",
		},
	}
	return app
}
