package main

import "executable/commands"

// List of all commands available
var allCommands = []commands.Command{
	commands.Sync,
	commands.Init,
	commands.Status,
	commands.Commit,
	commands.Switch,
	commands.Line,
	commands.Delete,
	commands.Patch,
	commands.Absorb,
	commands.Resolve,
}

// List of option tags
var allOptions = []commands.Option{
	{"help", "h", false},
}
