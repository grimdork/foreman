package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CmdInit options.
type CmdInit struct{}

// Run the init command.
func (cmd *CmdInit) Run(in []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	r := bufio.NewReader(os.Stdin)
	fmt.Printf("Server URL [%s]: ", cfg.ServerURL)
	t, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	t = strings.ToLower(strings.TrimSpace(t))
	if t != "" {
		cfg.ServerURL = t
	}
	fmt.Printf("Key ID [%s]: ", cfg.ID)
	t, err = r.ReadString('\n')
	if err != nil {
		return err
	}

	t = strings.ToLower(strings.TrimSpace(t))
	if t != "" {
		cfg.ID = t
	}
	fmt.Printf("Key [%s]: ", cfg.Key)
	t, err = r.ReadString('\n')
	if err != nil {
		return err
	}

	t = strings.ToLower(strings.TrimSpace(t))
	if t != "" {
		cfg.Key = t
	}
	return saveConfig(cfg)
}
