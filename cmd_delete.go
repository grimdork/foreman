package main

import (
	"fmt"

	"github.com/grimdork/foreman/api"
	"github.com/grimdork/opt"
)

// CmdDeleteCanary options.
type CmdDeleteCanary struct {
	opt.DefaultHelp
	Name  string `placeholder:"NAME" help:"Name of the canary to delete."`
	Force bool   `short:"f" help:"Force deletion of the canary."`
}

// Run the set command.
func (cmd *CmdDeleteCanary) Run(in []string) error {
	if cmd.Help || cmd.Name == "" {
		return opt.ErrUsage
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if !cmd.Force {
		if !confirm() {
			return nil
		}
	}

	req := map[string]string{
		"name": cmd.Name,
	}

	_, err = cfg.Delete(api.EPCanary, req)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted canary %s\n", cmd.Name)
	return nil
}

// CmdDeleteScout options.
type CmdDeleteScout struct {
	opt.DefaultHelp
	Hostname string `placeholder:"HOSTNAME" help:"Hostname of the scout to delete."`
	Force    bool   `short:"f" help:"Force deletion of the scout."`
}

// Run the set command.
func (cmd *CmdDeleteScout) Run(in []string) error {
	if cmd.Help || cmd.Hostname == "" {
		return opt.ErrUsage
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if !cmd.Force {
		if !confirm() {
			return nil
		}
	}

	req := map[string]string{
		"hostname": cmd.Hostname,
	}
	_, err = cfg.Delete(api.EPScout, req)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted scout %s\n", cmd.Hostname)
	return nil
}

// CmdDeleteKey options.
type CmdDeleteKey struct {
	opt.DefaultHelp
	ID    string `placeholder:"ID" help:"ID of the key to delete."`
	Force bool   `short:"f" help:"Force deletion of the key."`
}

// Run the set command.
func (cmd *CmdDeleteKey) Run(in []string) error {
	if cmd.Help || cmd.ID == "" {
		return opt.ErrUsage
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if !cmd.Force {
		if !confirm() {
			return nil
		}
	}

	req := map[string]string{
		"keyid": cmd.ID,
	}
	_, err = cfg.Delete(api.EPKey, req)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted key %s\n", cmd.ID)
	return nil
}
