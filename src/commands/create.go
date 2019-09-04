package commands

import (
	"errors"
	"fmt"
	"github.com/libgit2/git2go"
	"gitwrapper"
)

func execCreate(repo *git.Repository, positionals []string, _ map[string]string) error {
	if repo != nil {
		return errors.New("There is already a repository in this directory.")
	}

	directory := "."
	if len(positionals) > 0 {
		directory = positionals[0]
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}

	repo, err := gitwrapper.Init(directory)
	if err != nil {
		return err
	}
	err = gitwrapper.Commit(repo, "Create Repository")
	if err != nil {
		return err
	}

	fmt.Println("Created Metro repo.")
	return nil
}

func printCreateHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro create [directory]")
}

var Init = Command{"create", "Create a repo", execCreate, printCreateHelp}
