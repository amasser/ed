package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"

	"github.com/chzyer/readline"
	log "github.com/sirupsen/logrus"
)

// Editor ...
type Editor interface {
	io.Writer
	io.ReaderFrom

	Stop()
	Run() error
	Regexp() *regexp.Regexp
	SetRegexp(re *regexp.Regexp)
	Clipboard() []string
	SetClipboard(lines []string)
	Filename() string
	SetFilename(filename string)
	SetMode(mode int)
	SetPrompt(prompt string)
	Handle(cmd string, handler Handler)
}

type editor struct {
	rl        *readline.Instance
	mode      int
	running   bool
	buffer    Buffer
	filename  string
	clipboard []string
	regexp    *regexp.Regexp
	handlers  map[string]Handler
}

func newEditor() (Editor, error) {
	// TODO: Use functional options pattern here
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		InterruptPrompt: ".",
		EOFPrompt:       "q",

		VimMode: true,
	})
	if err != nil {
		return nil, err
	}

	e := &editor{
		rl:       rl,
		mode:     modeCommand,
		buffer:   newBuffer(),
		handlers: make(map[string]Handler),
	}

	return e, nil
}

func (e *editor) Write(p []byte) (n int, err error) {
	return
}

func (e *editor) ReadFrom(r io.Reader) (n int64, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		e.buffer.Append(line)
		n += int64(len(line) + 1)
	}
	if err = scanner.Err(); err != nil {
		log.Errorf("error reading from reader: %s", err)
		return
	}
	fmt.Printf("%d\n", n)
	return
}

func (e *editor) Regexp() *regexp.Regexp {
	return e.regexp
}

func (e *editor) SetRegexp(re *regexp.Regexp) {
	e.regexp = re
}

func (e *editor) Clipboard() []string {
	return e.clipboard
}

func (e *editor) SetClipboard(lines []string) {
	e.clipboard = lines[:]
}

func (e *editor) Filename() string {
	return e.filename
}

func (e *editor) SetFilename(filename string) {
	e.filename = filename
}

func (e *editor) SetMode(mode int) {
	e.mode = mode
}

func (e *editor) SetPrompt(prompt string) {
	e.rl.SetPrompt(prompt)
}

func (e *editor) Handle(cmd string, handler Handler) {
	e.handlers[cmd] = handler
}

func (e *editor) Stop() {
	e.running = false
}

func (e *editor) Close() {
	e.rl.Close()
}

func (e *editor) Run() error {
	defer e.Close()

	e.running = true
	for e.running {
		line, err := e.rl.Readline()
		if err != nil { // io.EOF
			if err == readline.ErrInterrupt {
				e.mode = modeCommand
				e.rl.SetPrompt("> ")
				continue
			} else if err == io.EOF {
				e.Stop()
				continue
			} else {
				return err
			}
		}

		if e.mode == modeCommand {
			cmd, err := parseCommand(line)
			if err != nil {
				log.Errorf("error parsing command: %s", err)
				continue
			}

			if err := cmd.Validate(e.buffer); err != nil {
				log.Errorf("error validating command: %s", err)
				continue
			}

			handler, ok := e.handlers[cmd.Cmd()]
			if !ok {
				log.Errorf("error unknown command: %s", line)
				continue
			}

			if err := handler(e, e.buffer, cmd); err != nil {
				log.Errorf("error processing command %s: %s", cmd.String(), err)
			}
		} else {
			switch e.mode {
			case modeAppend:
				e.buffer.Append(line)
			case modeInsert:
				e.buffer.Insert(line)
			default:
				panic("unknown input mode")
			}
		}
	}

	return nil
}
