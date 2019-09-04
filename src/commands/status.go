package commands

import (
	"fmt"
	"github.com/libgit2/git2go"
)

func execStatus(_ *git.Repository, _ []string, _ map[string]string) error {
	return nil
}

func printStatusHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro status")
}

var Status = Command{"status", "Show the state of the repo", execStatus, printStatusHelp}
