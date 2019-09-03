package main

import "commands"

// List of all commands available
var allCommands = []commands.Command{
	commands.Diff,
	commands.Sync,
	commands.Init,
	commands.Status,
	commands.Checkpoint,
}

// List of option tags
var allOptions = []commands.Option{
	{"help", "h", false},
	{"force", "f", false},
	{"timeout", "t", true},
}
