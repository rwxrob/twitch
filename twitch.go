package twitch

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
	"github.com/rwxrob/to"
	yq "github.com/rwxrob/yq/pkg"
)

var Cmd = &Z.Cmd{

	Name:      `twitch`,
	Summary:   `collection of twitch helper commands`,
	Version:   `v0.3.1`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Commands:  []*Z.Cmd{help.Cmd, bot, chat},
}

func sendChat(msg string) error {
	return Z.Exec([]string{"chat", msg}...)
}

var chat = &Z.Cmd{
	Name:    `chat`,
	Summary: `sends all arguments as a single string to Twitch chat`,
	Call: func(_ *Z.Cmd, args ...string) error {
		if len(args) == 0 {
			term.REPL(
				func(a string) string { return "" },
				func(a string) string { sendChat(a); return "" },
			)
		}
		msg := strings.Join(args, " ")
		return sendChat(msg)
	},
}

var bot = &Z.Cmd{
	Name:     `bot`,
	Summary:  `bot-related commands`,
	Commands: []*Z.Cmd{help.Cmd, commands},
}

var commands = &Z.Cmd{
	Name:    `commands`,
	Summary: `update and list Twitch Streamlabs Cloudbot commands`,
	Aliases: []string{"c", "cmd"},
	Commands: []*Z.Cmd{
		help.Cmd, add, edit, list, remove, _file, sync, commit,
	},
}

var commit = &Z.Cmd{
	Name:    `commit`,
	Summary: `commit the commands.yaml file`,
	Call: func(x *Z.Cmd, args ...string) error {
		path := x.Caller.C("file")
		if path == "" {
			return x.Caller.MissingConfig("file")
		}
		path = filepath.Dir(path)
		x.Log("Changing to directory: %v", path)
		if err := os.Chdir(path); err != nil {
			return err
		}
		Z.Exec(
			"git", "commit", "commands.yaml", "-m", "Update twitch/commands.yaml",
		)
		Z.Exec("git", "push")
		return nil
	},
}

var add = &Z.Cmd{
	Name:    `add`,
	Summary: `add a command by name from file`,
	Usage:   `<name>`,
	MinArgs: 1,
	Aliases: []string{"a"},
	Call: func(x *Z.Cmd, args ...string) error {
		if err := chat.Call(x,
			[]string{"!addcommand", args[0], "some"}...); err != nil {
			return err
		}
		return sync.Call(x, args[0])
	},
}

var remove = &Z.Cmd{
	Name:    `remove`,
	Summary: `remove a command with !rmcommand`,
	Usage:   `<command>`,
	Aliases: []string{"rm"},
	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) < 1 {
			return x.UsageError()
		}
		if args[0][0] != '!' {
			args[0] = "!" + args[0]
		}
		return chat.Call(x, []string{"!rmcommand", args[0]}...)
	},
}

var edit = &Z.Cmd{
	Name:    `edit`,
	Summary: `edit a command with !editcommand`,
	Usage:   `<command> <msg>`,
	MinArgs: 1,
	Aliases: []string{"rm"},
	Call: func(x *Z.Cmd, args ...string) error {
		if args[0][0] != '!' {
			args[0] = "!" + args[0]
		}
		msg := strings.Join(args[1:], " ")
		return chat.Call(x, []string{"!editcommand", args[0], msg}...)
	},
}

var sync = &Z.Cmd{
	Name:    `sync`,
	Summary: `sync a command from YAML file to Twitch`,
	Usage:   `<command>`,
	MinArgs: 1,
	Call: func(x *Z.Cmd, args ...string) error {
		path := x.Caller.C("file")
		if path == "" {
			return x.Caller.MissingConfig("file")
		}
		msg, err := yq.EvaluateToString("."+args[0], path)
		if err != nil {
			return err
		}
		if len(msg) > 380 {
			return fmt.Errorf("Must be 380 bytes or less (currently %v)", len(msg))
		}
		x.Log("Message body length: %v", len(msg))
		return edit.Call(x, args[0], msg)
	},
}

var _file = &Z.Cmd{
	Name:    `file`,
	Params:  []string{"edit"},
	Summary: `print the full path to commands file from configuration`,
	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) > 0 && args[0] == "edit" {
			file.Edit(x.Caller.C("file"))
			return commit.Call(x, args...)
		}
		fmt.Print(x.Caller.C("file"))
		return nil
	},
}

var list = &Z.Cmd{
	Name:    `list`,
	Summary: `list existing commands from commands.yaml`,
	Aliases: []string{"l"},
	Call: func(x *Z.Cmd, _ ...string) error {
		path := x.Caller.C("file")
		if path == "" {
			return x.Caller.MissingConfig("file")
		}
		buf, err := yq.EvaluateToString("keys", path)
		if err != nil {
			return err
		}
		lines := to.Lines(buf)
		sort.Strings(lines)
		buf = strings.Join(lines, " !")
		buf = strings.Replace(buf, "- ", "", -1)
		fmt.Println("!" + buf)
		return nil
	},
}
