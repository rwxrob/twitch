package twitch

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
	"github.com/rwxrob/to"
	yq "github.com/rwxrob/yq/pkg"
)

var Cmd = &Z.Cmd{

	Name:      `twitch`,
	Summary:   `collection of twitch helper commands`,
	Version:   `v0.3.2`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Commands:  []*Z.Cmd{help.Cmd, conf.Cmd, bot, chat},
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
	Commands: []*Z.Cmd{help.Cmd, conf.Cmd, commands},
}

var commands = &Z.Cmd{
	Name:    `commands`,
	Summary: `update and list Twitch Streamlabs Cloudbot commands`,
	Aliases: []string{"c", "cmd"},
	Commands: []*Z.Cmd{
		help.Cmd, add, edit, list, remove, filecmd, sync, commit,
	},
}

var commit = &Z.Cmd{
	Name:     `commit`,
	Summary:  `commit the commands.yaml file`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, args ...string) error {
		path, err := x.Caller.C("file")
		if err != nil {
			return err
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
	Name:     `add`,
	Summary:  `add a command by name from file`,
	Usage:    `<name>`,
	MinArgs:  1,
	Commands: []*Z.Cmd{help.Cmd},
	Aliases:  []string{"a"},
	Call: func(x *Z.Cmd, args ...string) error {
		if err := chat.Call(x,
			[]string{"!addcommand", args[0], "some"}...); err != nil {
			return err
		}
		return sync.Call(x, args[0])
	},
}

var remove = &Z.Cmd{
	Name:     `remove`,
	Summary:  `remove a command with !rmcommand`,
	Usage:    `<command>`,
	Commands: []*Z.Cmd{help.Cmd},
	Aliases:  []string{"rm"},
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
	Name:     `edit`,
	Summary:  `edit a command with !editcommand`,
	Usage:    `<command> <msg>`,
	Commands: []*Z.Cmd{help.Cmd},
	MinArgs:  1,
	Aliases:  []string{"rm"},
	Call: func(x *Z.Cmd, args ...string) error {
		if args[0][0] != '!' {
			args[0] = "!" + args[0]
		}
		msg := strings.Join(args[1:], " ")
		return chat.Call(x, []string{"!editcommand", args[0], msg}...)
	},
}

var sync = &Z.Cmd{
	Name:     `sync`,
	Summary:  `sync a command from YAML file to Twitch`,
	Usage:    `<command>`,
	Commands: []*Z.Cmd{help.Cmd},
	MinArgs:  1,
	Call: func(x *Z.Cmd, args ...string) error {
		path, err := x.Caller.C("file")
		if err != nil {
			return err
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

var filecmd = &Z.Cmd{
	Name:     `file`,
	Summary:  `print the full path to commands file from configuration`,
	NoArgs:   true,
	Commands: []*Z.Cmd{help.Cmd, fileedit},
	Call: func(x *Z.Cmd, args ...string) error {
		path, err := x.Caller.C("file")
		if err != nil {
			return err
		}
		if !term.IsInteractive() {
			path = strings.TrimSpace(path)
		}
		fmt.Print(path)
		return nil
	},
}

var fileedit = &Z.Cmd{
	Name:     `edit`,
	Summary:  `edit bot commands file with configured editor`,
	NoArgs:   true,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		path, err := x.Caller.Caller.C("file")
		if err != nil {
			return err
		}
		return file.Edit(strings.TrimSpace(path))
	},
}

var list = &Z.Cmd{
	Name:     `list`,
	Summary:  `list existing commands from commands.yaml`,
	Aliases:  []string{"l"},
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		path, err := x.Caller.C("file")
		if err != nil {
			return err
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
