package main

import (
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	log "github.com/sirupsen/logrus"
)

type execResult struct {
	io.ReadCloser
	Status    int
	Output    []byte
	readIndex int64
}

func (res *execResult) Close() error {
	return nil
}

func (res *execResult) Read(p []byte) (n int, err error) {
	if res.readIndex >= int64(len(res.Output)) {
		err = io.EOF
		return
	}

	n = copy(p, res.Output[res.readIndex:])
	res.readIndex += int64(n)
	return
}

func execShell(dir, cmd string) (res *execResult, err error) {
	res = &execResult{}

	sh := exec.Command("/bin/sh", "-c", cmd)
	if dir != "" {
		sh.Dir = dir
	}

	res.Output, err = sh.CombinedOutput()
	if err != nil {
		log.WithError(err).
			WithField("cmd", cmd).
			Error("error executing command")

		// Shamelessly borrowed from https://github.com/prologic/je/blob/master/job.go#L247
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				res.Status = status.ExitStatus()
			}
		}
	}

	return
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func detectLexer(filename, source string) (lexer chroma.Lexer) {
	if filename != "" {
		lexer := lexers.Match(filename)
		if lexer != nil {
			return lexer
		}
	}
	return lexers.Analyse(source)
}

func highlightSource(w io.Writer, filename, source, formatter, style string) error {
	l := detectLexer(filename, source)
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Determine formatter.
	f := formatters.Get(formatter)
	if f == nil {
		f = formatters.Fallback
	}

	// Determine style.
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return f.Format(w, s, it)
}
