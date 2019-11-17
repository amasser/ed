package main

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	debug   bool
	version bool

	prompt string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [file]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVarP(&version, "version", "v", false, "display version information")
	flag.BoolVarP(&debug, "debug", "d", false, "enable debug logging")

	flag.StringVarP(&prompt, "prompt", "p", "> ", "prompt to use")
}

func main() {
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if version {
		fmt.Printf("ed version %s", FullVersion())
		os.Exit(0)
	}

	e, err := newEditor()
	if err != nil {
		log.Printf("error creating editor: %s", err)
		os.Exit(1)
	}

	e.Handle("", cmdMove)
	e.Handle("a", cmdAppend)
	e.Handle("c", cmdChange)
	e.Handle("d", cmdDelete)
	e.Handle("i", cmdInsert)
	e.Handle("n", cmdNumber)
	e.Handle("p", cmdPrint)
	e.Handle("q", cmdQuit)
	e.Handle("w", cmdWrite)

	if len(flag.Args()) == 1 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			log.WithError(err).Error("error opening file")
			os.Exit(1)
		}
		if _, err = io.Copy(e, f); err != nil {
			log.WithError(err).Error("error reading from file")
			os.Exit(1)
		}
		f.Close()
	}

	if err := e.Run(); err != nil {
		log.Printf("error running editor: %s", err)
		os.Exit(1)
	}
}
