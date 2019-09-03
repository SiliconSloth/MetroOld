package commands

import (
	"errors"
	"fmt"
	"github.com/libgit2/git2go"
	"time"
)

func execInit(_ *git.Repository, positionals []string, options map[string]string) {
	println("Initing!")
	repo, err := git.InitRepository("test/.git", false) // TODO change to just .git

	if err != nil {
		println(err.Error())
	}

	err = createInitialCommit(repo)

	if err != nil {
		println(err.Error())
	}
}

func createInitialCommit(repo *git.Repository) error {
	index, err := repo.Index()

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create initial checkpoint as repo does not exist:\n%s", err.Error()))
	}

	oid, err := index.WriteTree()

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create initial checkpoint as initial tree could not be written:\n%s", err.Error()))
	}

	tree, err := repo.LookupTree(oid)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create initial checkpoint as tree does not exist:\n%s", err.Error()))
	}

	author := git.Signature{
		"Bob",
		"test@email.com",
		time.Now(),
	} // TODO change

	_, err = repo.CreateCommit("HEAD", &author, &author, "Initial Commit", tree)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create initial checkpoint:\n%s", err.Error()))
	}

	return nil
}

func printInitHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro init")
}

var Init = Command{"init", "Test git init", execInit, printInitHelp}
