package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
)

func execCommit(repo *git.Repository, positionals []string, _ map[string]string) error {
	if len(positionals) < 1 {
		return errors.New("Message required.")
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}
	message := positionals[0]

	err := gitwrapper.AssertConflicts(repo)
	if err != nil {
		return err
	}

	err = gitwrapper.Commit(repo, message, "HEAD^{commit}")
	if err != nil {
		return err
	}

	fmt.Println("Saved commit to current branch.")
	return nil
}

func printCommitHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro commit <message>")
}

var Commit = Command{"commit", "Make a commit", execCommit, printCommitHelp}
