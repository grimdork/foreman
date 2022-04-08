package main

import (
	ll "github.com/grimdork/loglines"
)

func (srv *Server) startLogger() {
	srv.log = make(chan string, 100)
	srv.err = make(chan string, 100)
	srv.logquit = make(chan any)
	srv.Add(1)
	srv.log <- "Starting logger."
	go func() {
		for {
			select {
			case msg := <-srv.log:
				ll.Msg(msg)

			case msg := <-srv.err:
				ll.Err(msg)

			case <-srv.logquit:
				ll.Msg("Stopping logger.")
				srv.Done()
				close(srv.logquit)
				close(srv.err)
				close(srv.log)
				return
			}
		}
	}()
}

func (srv *Server) stopLogger() {
	srv.logquit <- true
}
