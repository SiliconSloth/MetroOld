package commands

import (
	"fmt"
	git "github.com/libgit2/git2go"
)

func execCheckpoint(repo *git.Repository, positionals []string, options map[string]string) {
}

func printCheckpointHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro checkpoint")
}

var Checkpoint = Command{"checkpoint", "Checkpoint Command", execCheckpoint, printCheckpointHelp}

