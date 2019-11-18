package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	syntax "github.com/alecthomas/chroma/quick"
	log "github.com/sirupsen/logrus"
)

func cmdAppend(e Editor, buf Buffer, cmd Command) error {
	e.SetMode(modeAppend)
	e.SetPrompt("")
	return nil
}

func cmdDelete(e Editor, buf Buffer, cmd Command) error {
	buf.Delete(cmd.Addr())
	return nil
}

func cmdChange(e Editor, buf Buffer, cmd Command) error {
	buf.Delete(cmd.Addr())
	e.SetMode(modeInsert)
	e.SetPrompt("")
	return nil
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
	lines := buf.Select(cmd.Addr(), false)
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
	fmt.Println(buf.Current(false))

	return nil
}

func cmdNumber(e Editor, buf Buffer, cmd Command) error {
	selection := buf.Select(cmd.Addr(), true)
	source := strings.Join(selection, "\n") + "\n"
	err := syntax.Highlight(os.Stdout, source, "go", "terminal16m", "vim")
	if err != nil {
		log.WithError(err).Error("error syntax highlighting selection")
		return err
	}
	return nil

}

func cmdPrint(e Editor, buf Buffer, cmd Command) error {
	selection := buf.Select(cmd.Addr(), false)
	source := strings.Join(selection, "\n") + "\n"
	err := syntax.Highlight(os.Stdout, source, "go", "terminal16m", "vim")
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

	if e.Filename() == "" {
		e.SetFilename(filename)
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Errorf("error opening file for reading: %s", err)
		return err
	}
	defer f.Close()

	n, err := io.Copy(buf, f)
	if err != nil {
		log.Errorf("rror reading from input file: %s", err)
		return err
	}

	fmt.Printf("%d\n", n)

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
	e.SetClipboard(buf.Select(cmd.Addr(), false))
	return nil
}
