package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"time"
)

func execCheckpoint(repo *git.Repository, positionals []string, options map[string]string) error {
	if len(positionals) < 1 {
		return errors.New("Message required.")
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}
	message := positionals[1]

	author := git.Signature{
		"Bob",
		"test@email.com",
		time.Now(),
	}

	tree, err := getTree(repo)
	if err != nil {
		return err
	}

	_, err = repo.CreateCommit("HEAD", &author, &author, message, tree)
	if err != nil {
		return err
	}

	return nil
}

func getTree(repo *git.Repository) (*git.Tree, error) {
	index, err := repo.Index()
	if err != nil {
		return nil, err
	}

	oid, err := index.WriteTree()
	if err != nil {
		return nil, err
	}

	tree, err := repo.LookupTree(oid)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func printCheckpointHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro checkpoint <message>")
}

var Checkpoint = Command{"checkpoint", "Checkpoint Command", execCheckpoint, printCheckpointHelp}
