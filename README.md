# Common Twitch Commands, a Bonzai Branch

![WIP](https://img.shields.io/badge/status-wip-red)
![Go
Version](https://img.shields.io/github/go-mod/go-version/rwxrob/twitch)
[![GoDoc](https://godoc.org/github.com/rwxrob/twitch?status.svg)](https://godoc.org/github.com/rwxrob/twitch)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

This command is constantly under revision while I work on it. Don't plan
on anything stable until it hits `v1.0`.

## Install

This command can be installed as a standalone program (less preferred)
or composed into a Bonzai command tree (more preferred to avoid conflict
with the official `twitch` executable maintained by the Twitch company)

Standalone

```
go install github.com/rwxrob/twitch/cmd/twitch@latest
```

Composed

```go
package cmds

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/twitch"
)

var Cmd = &bonzai.Cmd{
	Name:     `cmds`,
	Commands: []*bonzai.Cmd{help.Cmd, twitch.Cmd},
}
```

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C twitch twitch
```

If you don't have bash or tab completion check use the shortcut
commands instead.

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.
