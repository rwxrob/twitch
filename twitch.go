package twitch

import (
	"fmt"
	"log"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/inc/help"
)

var Cmd = &bonzai.Cmd{

	Name:      `twitch`,
	Summary:   `collection of twitch helper commands`,
	Version:   `v0.0.1`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Commands:  []*bonzai.Cmd{help.Cmd, Command, Chat},
}

var Chat = &bonzai.Cmd{
	Name:    `chat`,
	Summary: `sends all arguments as a single string to Twitch chat`,
	Call: func(x *bonzai.Cmd, args ...string) error {
		if len(args) < 1 {
			return fmt.Errorf("empty chat message, need some arguments")
		}
		msg := strings.Join(args, " ")
		// FIXME: don't depend on command line `chat` program
		return bonzai.Exec([]string{"chat", msg}...)
	},
}

var Command = &bonzai.Cmd{
	Name:     `command`,
	Summary:  `update and list Twitch Streamlabs Cloudbot commands`,
	Aliases:  []string{"c", "cmd"},
	Commands: []*bonzai.Cmd{help.Cmd, CommandAdd, CommandList},
}

var CommandAdd = &bonzai.Cmd{
	Name:    `add`,
	Summary: `add (or update) a command with !addcommand`,
	Aliases: []string{"a"},
	Call: func(x *bonzai.Cmd, args ...string) error {
		if len(args) < 2 {
			return x.UsageError()
		}
		msg := strings.Join(args[1:], " ")
		return Chat.Call(x, []string{"!addcommand", args[0], msg}...)
	},
}

var CommandList = &bonzai.Cmd{
	Name:    `list`,
	Summary: `list existing commands from commands.yaml`,
	Aliases: []string{"l"},
	Call: func(_ *bonzai.Cmd, _ ...string) error {
		log.Print("would lookup commands.yaml")
		return nil
	},
}
