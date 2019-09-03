package commands

import (
	"errors"
	"fmt"
	"github.com/libgit2/git2go"
	"time"
)

func execCreate(_ *git.Repository, positionals []string, options map[string]string) error {
	directory := "."
	if len(positionals) > 0 {
		directory = positionals[0]
	}
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}

	repo, err := git.InitRepository(directory+"/.git", false)
	if err != nil {
		return err
	}

	err = createInitialCommit(repo)
	if err != nil {
		return err
	}

	fmt.Println("Created Metro repo.")
	return nil
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

	_, err = repo.CreateCommit("HEAD", &author, &author, "Initial checkpoint", tree)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create initial checkpoint:\n%s", err.Error()))
	}

	return nil
}

func printCreateHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro create [directory]")
}

var Init = Command{"create", "Create a blank repo", execCreate, printCreateHelp}
