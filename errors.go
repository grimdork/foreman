package main

import "errors"

var (
	// ErrClientNoExist is returned when a client does not exist.
	ErrClientNoExist = errors.New("client does not exist")
	// ErrClientExists is returned when a client already exists.
	ErrClientExists = errors.New("client already exists")
)
