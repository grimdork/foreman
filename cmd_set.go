package main

import (
	"fmt"

	"github.com/grimdork/foreman/api"
	"github.com/grimdork/opt"
)

// CmdSetCanary options.
type CmdSetCanary struct {
	opt.DefaultHelp
	Name     string `placeholder:"NAME" help:"Name of the canary to create/modify."`
	Interval int    `short:"i" placeholder:"INTERVAL" help:"Interval to expect check-ins." default:"60"`
	Key      bool   `short:"k" placeholder:"KEY" help:"Generate a new key."`
}

// Run the set command.
func (cmd *CmdSetCanary) Run(in []string) error {
	if cmd.Help || cmd.Name == "" {
		return opt.ErrUsage
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	req := map[string]string{
		"name":     cmd.Name,
		"interval": fmt.Sprintf("%d", cmd.Interval),
	}

	_, err = cfg.Post(api.EPCanary, req)
	if err != nil {
		return err
	}

	if cmd.Key {
		key := RandString(40)
		req = map[string]string{
			"keyid": cmd.Name,
			"value": key,
		}

		_, err = cfg.Post(api.EPKey, req)
		if err != nil {
			return err
		}

		fmt.Printf("Canary %s created/modified with key %s\n", cmd.Name, key)
	} else {
		fmt.Printf("Canary %s created/modified.\n", cmd.Name)
	}
	return nil
}

// CmdSetScout options.
type CmdSetScout struct {
	opt.DefaultHelp
	Hostname string `placeholder:"HOSTNAME" help:"Hostname of the scout to create/modify."`
	Port     int    `short:"p" placeholder:"PORT" help:"Port to watch." default:"443"`
	Interval int    `short:"i" placeholder:"INTERVAL" help:"Interval of the checks." default:"60"`
}

// Run the set command.
func (cmd *CmdSetScout) Run(in []string) error {
	if cmd.Help || cmd.Hostname == "" {
		return opt.ErrUsage
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	req := map[string]string{
		"hostname": cmd.Hostname,
		"port":     fmt.Sprintf("%d", cmd.Port),
		"interval": fmt.Sprintf("%d", cmd.Interval),
	}

	_, err = cfg.Post(api.EPScout, req)
	if err != nil {
		return err
	}

	fmt.Printf("Scout %s created/modified.\n", cmd.Hostname)
	return nil
}

// CmdSetKey options.
type CmdSetKey struct {
	opt.DefaultHelp
	ID    string `placeholder:"ID" help:"ID of the key to set."`
	Value string `placeholder:"VALUE" help:"Value for the key."`
	Admin bool   `short:"a" help:"Make this an admin key."`
}

// Run the set command.
func (cmd *CmdSetKey) Run(in []string) error {
	if cmd.Help || cmd.ID == "" {
		return opt.ErrUsage
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	req := map[string]string{
		"keyid": cmd.ID,
		"value": cmd.Value,
		"admin": fmt.Sprintf("%t", cmd.Admin),
	}

	_, err = cfg.Post(api.EPKey, req)
	if err != nil {
		return err
	}

	return nil
}
