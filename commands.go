package main

import (
	"fmt"
	"io"
	"log"
	"os"
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

func cmdInsert(e Editor, buf Buffer, cmd Command) error {
	e.SetMode(modeInsert)
	e.SetPrompt("")
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
		log.Printf("error moving to line %d: %s", cmd.Addr().Start(), err)
		return err
	}
	fmt.Println(buf.Current(false))

	return nil
}

func cmdNumber(e Editor, buf Buffer, cmd Command) error {
	for _, line := range buf.Select(cmd.Addr(), true) {
		fmt.Println(line)
	}
	return nil
}

func cmdPrint(e Editor, buf Buffer, cmd Command) error {
	for _, line := range buf.Select(cmd.Addr(), false) {
		fmt.Println(line)
	}
	return nil
}

func cmdQuit(e Editor, buf Buffer, cmd Command) error {
	e.Stop()
	return nil
}

func cmdWrite(e Editor, buf Buffer, cmd Command) error {
	filename := cmd.Arg(0)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error opening file for writing: %s", err)
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, buf)
	if err != nil {
		log.Printf("rror writing to output file: %s", err)
		return err
	}
	return nil
}
