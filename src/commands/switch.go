package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
)

func execSwitch(repo *git.Repository, positionals []string, _ map[string]string) error {
	if len(positionals) < 1 {
		return errors.New("Line name required.")
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}
	name := positionals[0]

	return gitwrapper.CheckoutBranch(name, repo)
}

func printSwitchHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro switch <line>")
}

var Switch = Command{"switch", "Switch to a different line or branch", execSwitch, printSwitchHelp}
