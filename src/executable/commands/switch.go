package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"metro"
)

func execSwitch(repo *git.Repository, positionals []string, _ map[string]string) error {
	if len(positionals) < 1 {
		return errors.New("Branch name required.")
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}
	name := positionals[0]

	err := metro.SwitchBranch(name, repo)
	if err != nil {
		return err
	}

	fmt.Println("Switched to branch " + name + ".")
	return err
}

func printSwitchHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro switch <line>")
}

var Switch = Command{"switch", "Switch to a different line", execSwitch, printSwitchHelp}
