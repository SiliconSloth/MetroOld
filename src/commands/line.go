package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
	"strings"
)

func execLine(repo *git.Repository, positionals []string, _ map[string]string) error {
	if len(positionals) < 1 {
		return errors.New("Line name required.")
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}
	name := positionals[0]

	if strings.HasSuffix(name, "-wip") {
		return errors.New("Line name can't end in -wip.")
	}

	_, err := gitwrapper.CreateBranch(name, repo)
	if err != nil {
		return err
	}

	fmt.Println("Created line " + name + ".")
	return nil
}

func printLineHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro line <name>")
}

var Line = Command{"line", "Create a new line", execLine, printLineHelp}
