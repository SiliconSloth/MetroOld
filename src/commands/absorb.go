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

	err := gitwrapper.AssertMerging(repo)
	if err != nil {
		return err
	}

	err = gitwrapper.StartMerge(name, repo)
	if err != nil {
		return err
	}

	index, err := repo.Index()
	if err != nil {
		return err
	}

	if index.HasConflicts() {
		fmt.Println("Conflicts occurred, please resolve.")
	} else {
		// If no conflicts occurred make the merge commit right away.
		err = gitwrapper.MergeCommit(repo)
		if err != nil {
			return err
		}

		current, err := gitwrapper.CurrentBranchName(repo)
		if err != nil {
			return err
		}
		fmt.Println("Successfully absorbed " + name + " into " + current + ".")
	}

	return nil
}

func printAbsorbHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro absorb <other-branch>")
}

var Absorb = Command{"absorb", "StartMerge the changes in another branch into this one", execAbsorb, printAbsorbHelp}
