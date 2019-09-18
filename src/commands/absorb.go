package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
	"helper"
	"strings"
)

func execAbsorb(repo *git.Repository, positionals []string, options map[string]string) error {
	if len(positionals) < 1 {
		return errors.New("Branch/line name required.")
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}
	name := positionals[0]

	if strings.HasSuffix(name, helper.WipString) {
		return errors.New("Can't absorb WIP branch.")
	}

	err := gitwrapper.AssertConflicts(repo)
	if err != nil {
		return err
	}

	conflicts, err := gitwrapper.Merge(name, repo)
	if err != nil {
		return err
	}
	if conflicts {
		fmt.Println("Conflicts occurred, please resolve.")
	} else {
		current, err := gitwrapper.CurrentBranchName(repo)
		if err != nil {
			return err
		}
		fmt.Println(name, "successfully absorbed into", current+".")
	}

	return nil
}

func printAbsorbHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro absorb <other-branch>")
}

var Absorb = Command{"absorb", "Merge the changes in another branch into this one", execAbsorb, printAbsorbHelp}
