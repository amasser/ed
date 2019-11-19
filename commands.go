package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func cmdAppend(e Editor, buf Buffer, cmd Command) error {
	e.SetMode(modeAppend)
	e.SetPrompt("")
	return nil
}

func cmdChange(e Editor, buf Buffer, cmd Command) error {
	if cmd.Addr().IsUnspecified() && buf.Index() == buf.Size() {
		buf.Delete(cmd.Addr())
		e.SetMode(modeAppend)
	} else {
		buf.Delete(cmd.Addr())
		e.SetMode(modeInsert)
	}

	e.SetPrompt("")
	return nil
}

func cmdDelete(e Editor, buf Buffer, cmd Command) error {
	buf.Delete(cmd.Addr())
	return nil
}

func cmdEdit(e Editor, buf Buffer, cmd Command) error {
	buf.Clear()
	return cmdRead(e, buf, cmd)
}

func cmdFile(e Editor, buf Buffer, cmd Command) error {
	e.SetFilename(cmd.Arg(0))
	if e.Filename() != "" {
		fmt.Println(e.Filename())
	}
	return nil
}

func cmdIndex(e Editor, buf Buffer, cmd Command) error {
	fmt.Printf("%d\n", buf.Index())
	return nil
}

func cmdInsert(e Editor, buf Buffer, cmd Command) error {
	e.SetMode(modeInsert)
	e.SetPrompt("")
	return nil
}

func cmdJoin(e Editor, buf Buffer, cmd Command) error {
	lines := buf.Select(cmd.Addr())
	buf.Delete(cmd.Addr())
	buf.Append(strings.Join(lines, ""))
	return nil
}

func cmdMove(e Editor, buf Buffer, cmd Command) error {
	// Special case of an empty command which is also the move command
	// If for example we press ENTER (unspecified address) and the buffer is empty
	// then do nothing.
	if cmd.Addr().Start() == 0 && buf.Size() == 0 {
		return nil
	}

	err := buf.Move(cmd.Addr())
	if err != nil {
		log.Errorf("error moving to line %d: %s", cmd.Addr().Start(), err)
		return err
	}
	fmt.Println(buf.Current())

	return nil
}

func cmdNumber(e Editor, buf Buffer, cmd Command) error {
	selection := buf.Select(cmd.Addr())

	out := &bytes.Buffer{}
	source := strings.Join(selection, "\n") + "\n"

	err := highlightSource(out, e.Filename(), source, "terminal16m", "vim")
	if err != nil {
		log.WithError(err).Error("error syntax highlighting selection")
		return err
	}

	ln := cmd.Addr().Start()
	scanner := bufio.NewScanner(bytes.NewBuffer(out.Bytes()))
	for scanner.Scan() {
		if ln == buf.Index() {
			fmt.Printf("%4d*  %s\n", ln, scanner.Text())
		} else {
			fmt.Printf("%4d  %s\n", ln, scanner.Text())
		}
		ln++
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("error printing lines: %s", err)
		return err
	}

	return nil
}

func cmdPrint(e Editor, buf Buffer, cmd Command) error {
	selection := buf.Select(cmd.Addr())
	source := strings.Join(selection, "\n") + "\n"

	err := highlightSource(os.Stdout, e.Filename(), source, "terminal16m", "vim")
	if err != nil {
		log.WithError(err).Error("error syntax highlighting selection")
		return err
	}

	return nil
}

func cmdPut(e Editor, buf Buffer, cmd Command) error {
	for _, line := range e.Clipboard() {
		buf.Append(line)
	}
	return nil
}

func cmdQuit(e Editor, buf Buffer, cmd Command) error {
	e.Stop()
	return nil
}

func cmdRead(e Editor, buf Buffer, cmd Command) error {
	filename := cmd.Arg(0)
	if filename == "" {
		filename = e.Filename()
	}

	if filename == "" {
		err := errNoFileSpecified
		log.WithError(err).Error("error must specify a filename or set a default filename")
		return err
	}

	if e.Filename() == "" && !strings.HasPrefix(filename, "!") {
		e.SetFilename(filename)
	}

	var (
		r   io.ReadCloser
		err error
	)

	if strings.HasPrefix(filename, "!") {
		command := filename[1:]
		r, err = execShell("", command)
		if err != nil {
			log.Errorf("error running shell command %s: %s", command, err)
			return err
		}
	} else {
		r, err = os.Open(filename)
		if err != nil {
			log.Errorf("error opening file for reading: %s", err)
			return err
		}
	}
	defer r.Close()

	n, err := io.Copy(buf, r)

	if err != nil {
		log.Errorf("rror reading from input file: %s", err)
		return err
	}

	fmt.Printf("%d\n", n)

	return nil
}

func cmdShell(e Editor, buf Buffer, cmd Command) error {
	command := cmd.Arg(0)
	if command == "" {
		log.Error("error no command specified")
		return errNoCommandSpecified
	}

	res, err := execShell("", command)
	if err != nil {
		log.Errorf("error executing command %s: %s", command, err)
		return err
	}

	fmt.Println(string(res.Output))
	fmt.Println("!")

	return nil
}

func cmdSearch(e Editor, buf Buffer, cmd Command) error {
	expr := cmd.Arg(0)
	if expr != "" {
		re, err := regexp.Compile(expr)
		if err != nil {
			log.Errorf("error parsing expression: %s", err)
			return err
		}

		e.SetRegexp(re)
	}

	re := e.Regexp()

	if re == nil {
		log.Error("error no search expression specified or previously set")
		return errNoExpressionSpecified
	}

	ok := buf.Search(re)
	if ok {
		fmt.Println(buf.Current())
	}
	return nil
}

func cmdWrite(e Editor, buf Buffer, cmd Command) error {
	filename := cmd.Arg(0)
	if filename == "" {
		filename = e.Filename()
	}

	if filename == "" {
		err := errNoFileSpecified
		log.WithError(err).Error("error must specify a filename or set a default filename")
		return err
	}

	if e.Filename() == "" {
		e.SetFilename(filename)
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("error opening file for writing: %s", err)
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, buf)
	if err != nil {
		log.Errorf("rror writing to output file: %s", err)
		return err
	}

	fmt.Printf("%d\n", n)

	return nil
}

func cmdWriteQuit(e Editor, buf Buffer, cmd Command) error {
	if err := cmdWrite(e, buf, cmd); err != nil {
		return err
	}
	return cmdQuit(e, buf, cmd)
}

func cmdYank(e Editor, buf Buffer, cmd Command) error {
	e.SetClipboard(buf.Select(cmd.Addr()))
	return nil
}
