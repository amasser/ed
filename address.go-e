package main

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Address ...
type Address interface {
	fmt.Stringer

	IsUnspecified() bool
	Resolve(buf Buffer) error
	Start() int
	End() int
}

type address struct {
	start  string
	delim  string
	end    string
	_start int
	_end   int
}

func (a *address) IsUnspecified() bool {
	return a._start == 0 && a._end == 0 && a.delim == ""
}

func (a *address) String() string {
	return fmt.Sprintf("%d%s%d", a._start, a.delim, a._end)
}

func (a *address) Resolve(buf Buffer) error {
	if a.start == "" && a.end == "" {
		if a.delim == "," {
			a._start, a._end = 1, buf.Size()
		} else if a.delim == ";" {
			a._start, a._end = buf.Index(), buf.Size()
		} else {
		}
	} else {
		if a.start == "" {
		} else if a.start == "." {
			a._start = buf.Index()
		} else {
			n, err := strconv.Atoi(a.start)
			if err != nil {
				log.WithError(err).Error("error parsing start address")
				return err
			}
			if n > buf.Size() {
				return errAddressOutOfRange
			}
			a._start = n
		}

		if a.end == "" {
		} else if a.end == "$" {
			a._end = buf.Size()
		} else {
			n, err := strconv.Atoi(a.end)
			if err != nil {
				log.WithError(err).Error("error parsing end address")
				return err
			}
			if n > buf.Size() {
				return errAddressOutOfRange
			}
			a._end = n
		}
	}

	if a._start > 0 && a.delim == "" && a.end == "" {
		a._end = a._start
	}

	a.start, a.end = string(a._start), string(a._end)

	return nil
}

func (a *address) Start() int {
	return a._start
}

func (a *address) End() int {
	return a._end
}
