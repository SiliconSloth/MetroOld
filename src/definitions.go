package main

import "commands"

var allCommands = []commands.Command{
	commands.Diff,
	commands.Sync,
	commands.Init,
}

var allOptions = []commands.Option{
	{"help", "h", false},
	{"force", "f", false},
	{"timeout", "t", true},
}
