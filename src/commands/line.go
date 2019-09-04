package commands

import (
	"errors"
	"fmt"
git "github.com/libgit2/git2go"
	"gitwrapper"
)

func execLine(repo *git.Repository, positionals []string, options map[string]string) error {
	if len(positionals) < 1 {
		return errors.New("Branch name required.")
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}
	name := positionals[0]

	_, err := gitwrapper.CreateBranch(name, repo)
	if err != nil {
		return err
	}

	println("Created ")

	return nil
}

func printLineHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro line <name>")
}

var Line = Command{"line", "Line Command", execLine, printLineHelp}
