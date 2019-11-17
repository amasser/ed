package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	cmdRegex = regexp.MustCompile(
		`(?P<start>[0-9]+|\.)?(,(?P<end>[0-9]+|\$))?(?P<command>[a-zA-Z]*)(?P<arguments>.*)$`,
	)
)

// Command ...
type Command interface {
	fmt.Stringer

	Validate(bufsize int) error

	Addr() Address
	Arg(i int) string
	Cmd() string
}

type command struct {
	addr address
	cmd  string
	args []string
}

func (c command) String() string {
	args := strings.Join(c.args, " ")
	return fmt.Sprintf("%s%s %s", c.addr.String(), c.cmd, args)
}

func (c command) Validate(bufsize int) error {
	if !c.addr.IsValid(bufsize) {
		return errAddressOutOfRange
	}
	return nil
}

func (c command) Addr() Address {
	return c.addr
}

func (c command) Arg(i int) string {
	if i < len(c.args) {
		return c.args[i]
	}
	return ""
}

func (c command) Cmd() string {
	return c.cmd
}

func parseCommand(line string) (cmd command, err error) {
	if !cmdRegex.MatchString(line) {
		err = errInvalidCommand
		return
	}

	match := cmdRegex.FindStringSubmatch(line)
	result := make(map[string]string)
	for i, name := range cmdRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	var (
		start int
		end   int
	)

	if result["start"] == "." || result["start"] == "" {
		start = 0
	} else {
		if start, err = strconv.Atoi(result["start"]); err != nil {
			return
		}
	}

	if result["end"] == "" {
		end = 0
	} else if result["end"] == "$" {
		end = -1
	} else {
		if end, err = strconv.Atoi(result["end"]); err != nil {
			return
		}
	}

	addr := address{start, end}

	args := strings.Split(strings.TrimSpace(result["arguments"]), " ")

	cmd = command{addr, result["command"], args}

	return
}
