package commands

import (
	"fmt"
	git "github.com/libgit2/git2go"
	"time"
)

func execCheckpoint(repo *git.Repository, positionals []string, options map[string]string) {
	author := git.Signature{
		"Bob",
		"test@email.com",
		time.Now(),
	}

	tree, err := getTree(repo)

	if err != nil {
		println(err.Error())
		return
	}

	_, err = repo.CreateCommit("HEAD", &author, &author, "Test Commit", tree)

	if err != nil {
		println(err.Error())
		return
	}
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
	fmt.Println("Usage: metro checkpoint")
}

var Checkpoint = Command{"checkpoint", "Checkpoint Command", execCheckpoint, printCheckpointHelp}

