package main

import "commands"

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
}

// List of option tags
var allOptions = []commands.Option{
	{"help", "h", false},
	{"force", "f", false},
	{"timeout", "t", true},
}
