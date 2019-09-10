package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
)

func execDelete(repo *git.Repository, positionals []string, options map[string]string) error {
	if len(positionals) < 1 || (positionals[0] != "commit" && positionals[0] != "line") {
		return errors.New("Incorrect Paramater.")
	}
	if positionals[0] == "commit" {
		fmt.Println("Usage: metro delete commit <num>")
	}
	if positionals[0] == "line" {
		fmt.Println("Usage: metro delete line line-name")
	}
}

func printDeleteHelp(positionals []string, _ map[string]string) {
	if len(positionals) < 1 || (positionals[0] != "commit" && positionals[0] != "line") {
		fmt.Println("Usage: metro delete <commit/line>")
	}
	if positionals[0] == "commit" {
		fmt.Println("Usage: metro delete commit <num>")
	}
	if positionals[0] == "line" {
		fmt.Println("Usage: metro delete line line-name")
	}
}

var Delete = Command{"delete", "Deletes a commit or branch", execDelete, printDeleteHelp}
