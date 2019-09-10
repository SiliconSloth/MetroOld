package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
	"strconv"
)

func execDelete(repo *git.Repository, positionals []string, options map[string]string) error {
	if len(positionals) < 1 || (positionals[0] != "commit" && positionals[0] != "line") {
		return errors.New("Incorrect paramater.")
	}
	if positionals[0] == "commit" {
		if len(positionals) > 2 {
			return errors.New("Unexpected argument: " + positionals[2])
		}
		deletes := 1
		if len(positionals) == 2 {
			var err error
			deletes, err = strconv.Atoi(positionals[1])
			if err != nil { return err }
		}
		return gitwrapper.RevertCommit(repo, deletes, false)
	}
	if positionals[0] == "line" {
		if len(positionals) < 2 {
			return errors.New("Line name required.")
		}
		if len(positionals) > 2 {
			return errors.New("Unexpected argument: " + positionals[2])
		}
		name := positionals[1]
		current, err := gitwrapper.CurrentBranchName(repo)
		if err != nil { return err }
		if name == current {
			return errors.New("Can't delete current branch.")
		}

		return gitwrapper.DeleteBranch(name, repo)
	}
	return nil
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
