package main

import (
	"log"
	"os"
)

func main() {
	e, err := newEditor()
	if err != nil {
		log.Printf("error creating editor: %s", err)
		os.Exit(1)
	}

	e.Handle("", cmdMove)
	e.Handle("a", cmdAppend)
	e.Handle("i", cmdInsert)
	e.Handle("n", cmdNumber)
	e.Handle("p", cmdPrint)
	e.Handle("q", cmdQuit)
	e.Handle("w", cmdWrite)

	if err := e.Run(); err != nil {
		log.Printf("error running editor: %s", err)
		os.Exit(1)
	}
}
