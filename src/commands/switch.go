package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
)

func execSwitch(repo *git.Repository, positionals []string, _ map[string]string) error {
	if len(positionals) < 1 {
		return errors.New("line name required")
	}
	if len(positionals) > 1 {
		return errors.New("unexpected argument: " + positionals[1])
	}
	name := positionals[0]

	err := gitwrapper.WIPCommit(repo)
	if err != nil {return err}
	err = gitwrapper.CheckoutBranch(name, repo)
	if err != nil {return err}
	err = gitwrapper.WIPUncommit(repo)
	return err
}

func printSwitchHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro switch <line>")
}

var Switch = Command{"switch", "Switch to a different line or branch", execSwitch, printSwitchHelp}
