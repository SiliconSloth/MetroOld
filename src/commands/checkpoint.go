package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
)

func execCheckpoint(repo *git.Repository, positionals []string, _ map[string]string) error {
	if len(positionals) < 1 {
		return errors.New("Message required.")
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}
	message := positionals[0]

	err := gitwrapper.Commit(repo, message, "HEAD^{commit}")
	if err != nil {
		return err
	}

	fmt.Println("Saved checkpoint to current branch.")
	return nil
}

func printCheckpointHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro checkpoint <message>")
}

var Checkpoint = Command{"checkpoint", "Checkpoint Command", execCheckpoint, printCheckpointHelp}
