package main

import (
	"github.com/Urethramancer/daemon"
)

// CmdServe options.
type CmdServe struct{}

// Run the serve command.
func (cmd *CmdServe) Run(in []string) error {
	srv, err := NewServer()
	if err != nil {
		return err
	}

	err = srv.Start()
	if err != nil {
		return err
	}

	<-daemon.BreakChannel()
	srv.Stop()
	return nil
}
