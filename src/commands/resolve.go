package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
)

func execResolve(repo *git.Repository, positionals []string, options map[string]string) error {
	merging := gitwrapper.MergeOngoing(repo)
	if !merging {
		return errors.New("You can only resolve conflicts while absorbing.")
	}

	err := gitwrapper.MergeCommit(repo)
	if err != nil {
		return err
	}

	current, err := gitwrapper.CurrentBranchName(repo)
	if err != nil {
		return err
	}
	fmt.Println("Successfully absorbed into " + current + ".")

	return nil
}

func printResolveHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro resolve")
}

var Resolve = Command{"resolve", "Commit resolved conflicts after absorb", execResolve, printResolveHelp}
