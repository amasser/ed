package main

import (
	"fmt"
	"io"
)

// Buffer ...
type Buffer interface {
	io.Reader
	io.WriterTo

	Index() int
	Size() int

	Append(line string)
	Current(lineno bool) string
	Delete(addr Address)
	Insert(line string)
	Move(addr Address) error
	Select(addr Address, showlns bool) []string
}

type buffer struct {
	index int
	lines []string
}

func (b *buffer) Index() int {
	return b.index
}

func (b *buffer) Size() int {
	return len(b.lines)
}

func (b *buffer) Append(line string) {
	b.lines = append(b.lines[:b.index], append([]string{line}, b.lines[b.index:]...)...)
	b.index++
}

func (b *buffer) Current(lineno bool) string {
	if len(b.lines) == 0 {
		return ""
	}

	if lineno {
		return fmt.Sprintf("%d\t%s", b.index, b.lines[(b.index-1)])
	}
	return b.lines[(b.index - 1)]
}

func (b *buffer) Delete(addr Address) {
	var start, end int

	if addr.IsUnspecified() {
		start, end = b.index, b.index
		b.lines = append(b.lines[:(start-1)], b.lines[end:]...)
		b.index--
	} else {
		if addr.End() == 0 {
			start, end = addr.Start(), addr.Start()
		} else {
			start = addr.Start()
			if addr.End() == -1 {
				end = len(b.lines)
			} else {
				end = addr.End()
			}
		}

		b.lines = append(b.lines[:start], b.lines[(end+1):]...)

		if len(b.lines) == 0 {
			b.index = 0
		} else {
			b.index -= (end - (start - 1))
		}
	}
}

func (b *buffer) Insert(line string) {
	b.lines = append(b.lines[:(b.index-1)], append([]string{line}, b.lines[(b.index-1):]...)...)
}

func (b *buffer) Move(addr Address) error {
	var n int

	if addr.Start() == 0 {
		n = b.index
	} else {
		n = addr.Start()
	}

	// cmdMove won't call us anyway as it checks for an empty buffer and the
	// special case 0th index (also representing an empty buffer)
	if n == 0 && len(b.lines) == 0 {
		return nil
	}
	if n < 1 || n > len(b.lines) {
		return errAddressOutOfRange
	}
	b.index = n
	return nil
}

func (b *buffer) Select(addr Address, showlns bool) []string {
	if len(b.lines) == 0 {
		return nil
	}

	if addr.IsUnspecified() {
		if showlns {
			return []string{fmt.Sprintf("%d\t%s", b.index, b.lines[(b.index-1)])}
		}
		return []string{b.lines[(b.index - 1)]}
	}

	var (
		start, end int
		lines      []string
	)

	// XXX: Generalize this
	if addr.End() == 0 {
		start, end = addr.Start(), addr.Start()
	} else {
		start = addr.Start()
		if addr.End() == -1 {
			end = len(b.lines)
		} else {
			end = addr.End()
		}
	}

	for i := start; i <= end; i++ {
		if showlns {
			lines = append(lines, fmt.Sprintf("%d\t%s", i, b.lines[(i-1)]))
		} else {
			lines = append(lines, b.lines[(i-1)])
		}
	}

	return lines
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
