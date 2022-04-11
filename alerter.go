package main

import (
	"github.com/grimdork/foreman/api"
	ll "github.com/grimdork/loglines"
)

func (srv *Server) startAlerter() {
	srv.logquit = make(chan any)
	srv.alerts = make(chan api.Message, 100)
	srv.Add(1)
	ll.Msg("Starting alerter.")
	go func() {
		for {
			select {
			case msg := <-srv.alerts:
				if msg.OldStatus == msg.NewStatus {
					continue
				}

				switch msg.Type {
				case api.FailNoTLS:
					ll.Msg("%s: Not a secure connection", msg.Name)
				}

			case <-srv.logquit:
				ll.Msg("Stopping alerter.")
				srv.Done()
				close(srv.logquit)
				close(srv.alerts)
				return
			}
		}
	}()
}

func (srv *Server) stopLogger() {
	srv.logquit <- true
}
