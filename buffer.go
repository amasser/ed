package main

import (
	"fmt"
	"io"
)

// Buffer ...
type Buffer interface {
	io.Reader
	io.WriterTo

	Append(line string)
	Insert(line string)
	Move(n int) error
	Current(lineno bool) string
}

type buffer struct {
	index int
	lines []string
}

func (b *buffer) Append(line string) {
	b.lines = append(b.lines[:b.index], append([]string{line}, b.lines[b.index:]...)...)
	b.index++
}

func (b *buffer) Insert(line string) {
	b.lines = append(b.lines[:(b.index-1)], append([]string{line}, b.lines[(b.index-1):]...)...)
}

func (b *buffer) Move(n int) error {
	if n < 1 || n > len(b.lines) {
		return errAddressOutOfRange
	}
	b.index = n
	return nil
}

func (b *buffer) Current(lineno bool) string {
	if lineno {
		return fmt.Sprintf("%d\t%s", b.index, b.lines[(b.index-1)])
	}
	return b.lines[(b.index - 1)]
}

func (b *buffer) Read(p []byte) (n int, err error) {
	return
}

func (b *buffer) WriteTo(w io.Writer) (n int64, err error) {
	for _, line := range b.lines {
		if _, err = w.Write([]byte(line)); err != nil {
			return
		}
		if _, err = w.Write([]byte("\n")); err != nil {
			return
		}
		n += int64(len(line)) + 1
	}
	return
}

func newBuffer() Buffer {
	return &buffer{lines: make([]string, 0)}
}
