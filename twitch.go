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
	Commands:  []*bonzai.Cmd{help.Cmd, bot, chat},
}

var chat = &bonzai.Cmd{
	Name:    `chat`,
	Summary: `sends all arguments as a single string to Twitch chat`,
	Call: func(x *bonzai.Cmd, args ...string) error {
		// TODO read from stdin if no arguments
		if len(args) < 1 {
			return fmt.Errorf("empty chat message, need some arguments")
		}
		msg := strings.Join(args, " ")
		// FIXME: don't depend on command line `chat` program
		return bonzai.Exec([]string{"chat", msg}...)
	},
}

var bot = &bonzai.Cmd{
	Name:     `bot`,
	Summary:  `bot-related commands`,
	Commands: []*bonzai.Cmd{help.Cmd, commands},
}

var commands = &bonzai.Cmd{
	Name:     `commands`,
	Summary:  `update and list Twitch Streamlabs Cloudbot commands`,
	Aliases:  []string{"c", "cmd"},
	Commands: []*bonzai.Cmd{help.Cmd, add, edit, list, remove},
}

var add = &bonzai.Cmd{
	Name:    `add`,
	Summary: `add (or update) a command with !addcommand`,
	Usage:   `<command> <body>`,
	Aliases: []string{"a"},
	Call: func(x *bonzai.Cmd, args ...string) error {
		if len(args) < 2 {
			return x.UsageError()
		}
		msg := strings.Join(args[1:], " ")
		return chat.Call(x, []string{"!addcommand", args[0], msg}...)
	},
}

var remove = &bonzai.Cmd{
	Name:    `remove`,
	Summary: `remove a command with !rmcommand`,
	Usage:   `<command>`,
	Aliases: []string{"rm"},
	Call: func(x *bonzai.Cmd, args ...string) error {
		if len(args) < 1 {
			return x.UsageError()
		}
		if args[0][0] != '!' {
			args[0] = "!" + args[0]
		}
		return chat.Call(x, []string{"!rmcommand", args[0]}...)
	},
}

var edit = &bonzai.Cmd{
	Name:    `edit`,
	Summary: `edit a command with !editcommand`,
	Usage:   `<command> <msg>`,
	Aliases: []string{"rm"},
	Call: func(x *bonzai.Cmd, args ...string) error {
		if len(args) < 1 {
			return x.UsageError()
		}
		if args[0][0] != '!' {
			args[0] = "!" + args[0]
		}
		msg := strings.Join(args[1:], " ")
		return chat.Call(x, []string{"!editcommand", args[0], msg}...)
	},
}

var list = &bonzai.Cmd{
	Name:    `list`,
	Summary: `list existing commands from commands.yaml`,
	Aliases: []string{"l"},
	Call: func(x *bonzai.Cmd, _ ...string) error {
		path := x.Caller.Q("file")
		if path == "" {
			return x.Caller.MissingConfig("file")
		}
		log.Printf("would lookup %v", path)
		// TODO look into yq for this lookup
		return nil
	},
}
