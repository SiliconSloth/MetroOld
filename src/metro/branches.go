package metro

import (
	"errors"
	git "github.com/libgit2/git2go"
	"strings"
)

// Create a new branch from the current head with the specified name.
// Returns the branch
func CreateBranch(name string, repo *git.Repository) (*git.Branch, error) {
	commit, err := GetCommit("HEAD", repo)
	if err != nil {
		return nil, err
	}

	branch, err := repo.CreateBranch(name, commit, false)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

func SwitchBranch(name string, repo *git.Repository) error {
	if strings.HasSuffix(name, WipString) {
		return errors.New("Can't switch to wip line.")
	}
	if !CommitExists(name, repo) {
		return errors.New("No branch called " + name + ".")
	}

	err := SaveWIP(repo)
	if err != nil {
		return err
	}
	err = checkoutBranch(name, repo)
	if err != nil {
		return err
	}
	err = RestoreWIP(repo)
	if err != nil {
		return err
	}

	return nil
}

// Checks out the given branch by name
// name - Plain Text branch name (e.g. 'master')
// repo - Repo to checkout from
func checkoutBranch(name string, repo *git.Repository) error {
	err := checkout(name, repo)
	if err != nil {
		return err
	}

	err = moveHead(name, repo)
	return err
}

// Moves the head to the given branch
// Files are not changed
func moveHead(name string, repo *git.Repository) error {
	branch, err := repo.LookupBranch(name, git.BranchLocal)
	if err != nil {
		return err
	}

	err = repo.SetHead(branch.Reference.Name())
	return err
}

func DeleteBranch(name string, repo *git.Repository) error {
	branch, err := repo.LookupBranch(name, git.BranchLocal)
	if err != nil {
		return err
	}
	err = branch.Delete()
	if err != nil {
		return err
	}

	return nil
}

func CurrentBranchName(repo *git.Repository) (string, error) {
	iterator, err := repo.NewBranchIterator(git.BranchLocal)
	if err != nil {
		return "", err
	}
	for branch, _, err := iterator.Next(); err == nil; branch, _, err = iterator.Next() {
		head, err := branch.IsHead()
		if err != nil {
			return "", err
		}
		if head {
			return branch.Name()
		}
	}
	return "", errors.New("Could not find current line.")
}
