package main

import (
	"bufio"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

// Editor ...
type Editor interface {
	io.Writer
	io.ReaderFrom

	Stop()
	Run() error
	SetMode(mode int)
	SetPrompt(prompt string)
	Handle(cmd string, handler Handler)
}

type editor struct {
	rl       *readline.Instance
	mode     int
	running  bool
	buffer   Buffer
	handlers map[string]Handler
}

func newEditor() (Editor, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:  "> ",
		VimMode: true,
	})
	if err != nil {
		return nil, err
	}

	log.SetOutput(rl.Stderr())

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
		log.Printf("error reading from reader: %s", err)
		return
	}
	return
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
			if err == io.EOF {
				return nil
			}
			return err
		}

		if e.mode == modeCommand {
			cmd, err := parseCommand(line)
			if err != nil {
				log.Printf("error parsing command: %s", err)
				continue
			}

			if err := cmd.Validate(e.buffer.Size()); err != nil {
				log.Printf("error validating command: %s", err)
				continue
			}

			handler, ok := e.handlers[cmd.Cmd()]
			if !ok {
				log.Printf("error unknown command: %s", line)
				continue
			}

			if err := handler(e, e.buffer, cmd); err != nil {
				log.Printf("error processing command %s: %s", cmd.String(), err)
			}
		} else {
			if strings.TrimSpace(line) == "." {
				e.mode = modeCommand
				e.rl.SetPrompt("> ")
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
	}

	return nil
}
