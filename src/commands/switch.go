package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
	"helper"
	"strings"
)

func execSwitch(repo *git.Repository, positionals []string, _ map[string]string) error {
	if len(positionals) < 1 {
		return errors.New("Line name required.")
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}
	name := positionals[0]

	if strings.HasSuffix(name, helper.WipString) {
		return errors.New("Can't switch to wip line.")
	}
	if !gitwrapper.BranchExists(name, repo) {
		return errors.New("No line called " + name + ".")
	}

	err := gitwrapper.WIPCommit(repo)
	if err != nil {
		return err
	}
	err = gitwrapper.CheckoutBranch(name, repo)
	if err != nil {
		return err
	}
	err = gitwrapper.WIPUncommit(repo)
	return err
}

func printSwitchHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro switch <line>")
}

var Switch = Command{"switch", "Switch to a different line", execSwitch, printSwitchHelp}
