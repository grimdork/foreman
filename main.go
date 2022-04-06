package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/grimdork/opt"
)

var o struct {
	opt.DefaultHelp
	Serve CmdServe `command:"serve" help:"Start the server."`
	Init  CmdInit  `command:"init" help:"Initialize the configuration."`

	ListCanary   CmdListCanaries `command:"listcanaries" help:"List all passive clients." aliases:"lsc" group:"Canaries"`
	SetCanary    CmdSetCanary    `command:"setcanary" help:"Set a passive client." aliases:"sc" group:"Canaries"`
	DeleteCanary CmdDeleteCanary `command:"deletecanary" help:"Delete a passive client." aliases:"dc" group:"Canaries"`

	ListScouts  CmdListScouts  `command:"listscouts" help:"List all scouts (active watcher)." aliases:"lss" group:"Scouts"`
	SetScout    CmdSetScout    `command:"setscout" help:"Creates or modifies a scout." aliases:"ss" group:"Scouts"`
	DeleteScout CmdDeleteScout `command:"deletescout" help:"Deletes a scout." aliases:"ds" group:"Scouts"`

	ListKey   CmdListKeys  `command:"listkeys" help:"List all keys." aliases:"lsk" group:"Keys"`
	SetKey    CmdSetKey    `command:"setkey" help:"Set a key." aliases:"sk" group:"Keys"`
	DeleteKey CmdDeleteKey `command:"deletekey" help:"Delete a key." aliases:"dk" group:"Keys"`
}

func main() {
	err := opt.HandleCommands(&o)
	if err != nil {
		if err == opt.ErrNoCommand {
			fmt.Printf("Unknown command or no command specified.\n}")
			os.Exit(1)
		}

		panic(err)
	}
}

func parse(x any) error {
	a := opt.Parse(x)
	v := reflect.ValueOf(x).Elem().FieldByName("Help")
	if v.Bool() {
		a.Usage()
		return nil
	}

	err := a.RunCommand(false)
	if err != nil {
		if err == opt.ErrNoCommand {
			a.Usage()
			return nil
		}

		return err
	}

	return nil
}
