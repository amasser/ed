package main

import "errors"

var (
	errInvalidCommand    = errors.New("error: invalid command")
	errAddressOutOfRange = errors.New("error: address out of range")
)
