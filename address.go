package main

import (
	"fmt"
)

// Address ...
type Address interface {
	fmt.Stringer

	Start() int
	End() int
}

type address struct {
	start int
	end   int
}

func (a address) String() string {
	var (
		start string
		end   string
	)

	if a.start > 0 {
		start = fmt.Sprintf("%d", a.start)
	}

	if a.end == -1 {
		end = "$"
	} else if a.end > 0 {
		end = fmt.Sprintf("%d", a.end)
	}

	return fmt.Sprintf("%s%s", start, end)
}

func (a address) Start() int {
	return a.start
}

func (a address) End() int {
	return a.end
}
