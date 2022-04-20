package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/grimdork/foreman/api"
)

// CmdListCanaries options.
type CmdListCanaries struct{}

// Run the list command.
func (cmd *CmdListCanaries) Run(in []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	data, err := cfg.Get(api.EPCanaries, api.Request{})
	if err != nil {
		return err
	}

	if len(data) == 0 {
		fmt.Println("No canaries.")
		return nil
	}

	var list api.CanaryList
	err = json.Unmarshal(data, &list)
	if err != nil {
		return err
	}

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprint(w, "Host\tInterval\tStatus\tLast check\tFirst failure\tAssignee\tAssigned\n")
	for _, c := range list.Canaries {
		if c.Assignee == "" {
			c.Assignee = "Nobody"
		}
		fail := "N/A"
		if c.Status != api.StatusOK && c.Status != api.StatusWaiting {
			fail = c.Failed.Local().Format(time.RFC822)
		}
		fmt.Fprintf(w, "%s\t%ds\t%s\t%s\t%s\t%s\t%s\n",
			c.Hostname, c.Interval, api.StatusString(c.Status),
			c.LastCheck.Local().Format(time.RFC822), fail, c.Assignee, c.Assigned.String(),
		)
	}

	w.Flush()
	return nil
}

// CmdListScouts options.
type CmdListScouts struct{}

// Run the list command.
func (cmd *CmdListScouts) Run(in []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	data, err := cfg.Get(api.EPScouts, api.Request{})
	if err != nil {
		return err
	}

	var list api.ScoutList
	err = json.Unmarshal(data, &list)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		fmt.Println("No scouts.")
		return nil
	}

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprint(w, "Host\tPort\tInterval\tStatus\tLast check\tFirst failure\tAck\tAssignee\n")
	for _, s := range list.Scouts {
		if s.Assignee == "" {
			s.Assignee = "Nobody"
		}
		fail := "N/A"
		if s.Status != api.StatusOK && s.Status != api.StatusWaiting {
			fail = s.Failed.Local().Format(time.RFC822)
		}
		fmt.Fprintf(w, "%s\t%d\t%ds\t%s\t%s\t%s\t%s\t%s\n",
			s.Hostname, s.Port, s.Interval, api.StatusString(s.Status),
			s.LastCheck.Local().Format(time.RFC822), fail, s.Assignee, s.Assigned.String(),
		)
	}
	w.Flush()
	return nil
}

// CmdListKeys options.
type CmdListKeys struct{}

// Run the list command.
func (cmd *CmdListKeys) Run(in []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	data, err := cfg.Get(api.EPKeys, api.Request{})
	if err != nil {
		return err
	}

	if len(data) == 0 {
		fmt.Println("No keys.")
		return nil
	}

	var list api.KeyList
	err = json.Unmarshal(data, &list)
	if err != nil {
		return err
	}

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprint(w, "ID\tKey\tAdmin\n")
	for _, k := range list.Keys {
		fmt.Fprintf(w, "%s\t%s\t%t\n", k.ID, k.Key, k.Admin)
	}

	w.Flush()
	return nil
}
