package main

import (
	"fmt"
	"io"
	"log"
)

// Buffer ...
type Buffer interface {
	io.Reader
	io.WriterTo

	Append(line string)
	Current(lineno bool) string
	Delete(start, end int)
	Insert(line string)
	Move(n int) error
	Select(addr Address, showlns bool) []string
	Size() int
}

type buffer struct {
	index int
	lines []string
}

func (b *buffer) Append(line string) {
	b.lines = append(b.lines[:b.index], append([]string{line}, b.lines[b.index:]...)...)
	b.index++
}

func (b *buffer) Current(lineno bool) string {
	if lineno {
		return fmt.Sprintf("%d\t%s", b.index, b.lines[(b.index-1)])
	}
	return b.lines[(b.index - 1)]
}

func (b *buffer) Delete(start, end int) {
	if start == 0 && end == -1 {
		start, end = b.index, b.index
		b.lines = append(b.lines[:(start-1)], b.lines[end:]...)
		b.index--
	} else {
		log.Printf("deleting from %d to %d", start, end)
	}
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

func (b *buffer) Select(addr Address, showlns bool) []string {
	if addr.IsUnspecified() {
		if showlns {
			return []string{fmt.Sprintf("%d\t%s", b.index, b.lines[(b.index-1)])}
		}
		return []string{b.lines[(b.index - 1)]}
	}

	var lines []string

	for i := addr.Start(); i <= addr.End(); i++ {
		if showlns {
			lines = append(lines, fmt.Sprintf("%d\t%s", i, b.lines[(i-1)]))
		} else {
			lines = append(lines, b.lines[(i-1)])
		}
	}

	return lines
}

func (b *buffer) Size() int {
	return len(b.lines)
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
