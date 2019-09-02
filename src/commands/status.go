package commands

import (
	"fmt"
	"github.com/libgit2/git2go"
)

func execStatus(repo *git.Repository, positionals []string, options map[string]string) {
}

func printStatusHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro status")
}

var Status = Command{"status", "Status Command", execStatus, printStatusHelp}

