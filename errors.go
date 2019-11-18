package main

import "errors"

var (
	errInvalidCommand        = errors.New("error: invalid command")
	errAddressOutOfRange     = errors.New("error: address out of range")
	errNoFileSpecified       = errors.New("error: no filename specified")
	errNoExpressionSpecified = errors.New("error: no expression specified")
)
