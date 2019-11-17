package main

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type execResult struct {
	Status int
	Output []byte
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
