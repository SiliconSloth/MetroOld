package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
)

func execPatch(repo *git.Repository, positionals []string, options map[string]string) error {
	// Uses existing message as default
	commit, err := gitwrapper.GetLastCommit(repo)
	if err != nil {return err}
	message := commit.Message()

	if len(positionals) == 1 {
		// Overrides message
		message = positionals[0]
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}

	err = gitwrapper.RevertLastCommit(repo, false)
	if err != nil {return err}
	err = gitwrapper.Commit(repo, message, "HEAD^{commit}")
	if err != nil {return err}

	fmt.Println("Patched Commit.")
	return nil
}

func printPatchHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro patch <message>")
}

var Patch = Command{"patch", "Will patch the last commit with the current work", execPatch, printPatchHelp}
