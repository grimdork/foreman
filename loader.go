package main

import (
	"github.com/grimdork/foreman/clients"
	ll "github.com/grimdork/loglines"
)

// LoadScouts loads all scouts from the database and starts watching.
func (srv *Server) LoadScouts() error {
	list, err := srv.GetScouts()
	if err != nil {
		return err
	}

	for _, e := range list.Scouts {
		scout := clients.NewScout(e.Hostname, srv.alerts, srv.capool)
		scout.Port = e.Port
		scout.Interval = e.Interval
		scout.LastCheck = e.LastCheck
		scout.Failed = e.Failed
		scout.Status = e.Status
		scout.Assignee = e.Assignee
		scout.Assigned = e.Assigned
		srv.scouts[e.Hostname] = scout
		ll.Msg("Starting scout for %s", scout.Hostname)
		scout.Start()
	}

	return nil
}

// LoadCanaries loads all canaries from the database.
func (srv *Server) LoadCanaries() error {
	list, err := srv.GetCanaries()
	if err != nil {
		return err
	}

	for _, e := range list.Canaries {
		canary := clients.NewCanary(e.Hostname, srv.alerts)
		canary.Key = e.Key
		canary.Interval = e.Interval
		canary.LastCheck = e.LastCheck
		canary.Failed = e.Failed
		canary.Status = e.Status
		canary.Assignee = e.Assignee
		canary.Assigned = e.Assigned
		srv.canaries[e.Hostname] = canary
	}
	return nil
}
