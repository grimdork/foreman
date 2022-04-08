package main

import "github.com/grimdork/foreman/clients"

// LoadScouts loads all scouts from the database and starts watching.
func (srv *Server) LoadScouts() error {
	list, err := srv.GetScouts()
	if err != nil {
		return err
	}

	for _, e := range list.Scouts {
		scout := clients.NewScout(e.Hostname, srv.log)
		scout.Port = e.Port
		scout.Interval = e.Interval
		scout.LastCheck = e.LastCheck
		scout.Failed = e.Failed
		scout.Status = e.Status
		scout.Acknowledgement = e.Acknowledgement
		scout.Assignee = e.Assignee
		srv.scouts[e.Hostname] = scout
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
		canary := clients.NewCanary(e.Hostname, srv.log)
		canary.Port = e.Port
	}
	return nil
}
