package plugcmds

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/gobuffalo/buffalo-plugins/plugins"
	"github.com/gobuffalo/events"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewAvailable returns a fully formed Available type
func NewAvailable() *Available {
	return &Available{
		plugs: map[string]plug{},
		moot:  &sync.RWMutex{},
	}
}

// Available used to manage all of the available commands
// for the plugin
type Available struct {
	plugs map[string]plug
	moot  *sync.RWMutex
}

type plug struct {
	BuffaloCommand string
	Cmd            *cobra.Command
	Plugin         plugins.Command
}

func (p plug) String() string {
	b, _ := json.Marshal(p.Plugin)
	return string(b)
}

// Cmd returns the "available" command
func (a Available) Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "available",
		Short: "a list of available buffalo plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Encode(os.Stdout)
		},
	}
}

// Commands returns all of the commands that are available
func (a Available) Commands() []*cobra.Command {
	cmds := []*cobra.Command{a.Cmd()}
	a.moot.RLock()
	for _, p := range a.plugs {
		cmds = append(cmds, p.Cmd)
	}
	a.moot.RUnlock()
	return cmds
}

// Add a new command to this list of available ones.
// The bufCmd should corresponding buffalo command that
// command should live below.
//
// Special "commands":
//	"root" - is the `buffalo` command
//	"events" - listens for emitted events
func (a *Available) Add(bufCmd string, cmd *cobra.Command) error {
	if len(cmd.Aliases) == 0 {
		cmd.Aliases = []string{}
	}
	p := plug{
		BuffaloCommand: bufCmd,
		Cmd:            cmd,
		Plugin: plugins.Command{
			Name:           cmd.Use,
			BuffaloCommand: bufCmd,
			Description:    cmd.Short,
			Aliases:        cmd.Aliases,
			UseCommand:     cmd.Use,
		},
	}
	a.moot.Lock()
	a.plugs[p.String()] = p
	a.moot.Unlock()
	return nil
}

// Mount all of the commands that are available
// on to the other command. This is the recommended
// approach for using Available.
//	a.Mount(rootCmd)
func (a Available) Mount(cmd *cobra.Command) {
	// mount all the cmds on to the cobra command
	cmd.AddCommand(a.Cmd())
	a.moot.RLock()
	for _, p := range a.plugs {
		cmd.AddCommand(p.Cmd)
	}
	a.moot.RUnlock()
}

// Encode into the required Buffalo plugins available
// formate
func (a *Available) Encode(w io.Writer) error {
	var plugs plugins.Commands
	a.moot.RLock()
	for _, p := range a.plugs {
		plugs = append(plugs, p.Plugin)
	}
	a.moot.RUnlock()
	return json.NewEncoder(w).Encode(plugs)
}

// Listen adds a command for github.com/gobuffalo/events.
func (a *Available) Listen(fn func(e events.Event) error) error {
	listenCmd := &cobra.Command{
		Use:   "listen",
		Short: "listens to github.com/gobuffalo/events",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("must pass a payload")
			}

			e := events.Event{}
			err := json.Unmarshal([]byte(args[0]), &e)
			if err != nil {
				return errors.WithStack(err)
			}

			return fn(e)
		},
	}
	return a.Add("events", listenCmd)
}
