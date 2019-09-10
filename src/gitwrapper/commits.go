package gitwrapper

import (
	"errors"
	git "github.com/libgit2/git2go"
	"time"
)

// Commit all files in the repo directory (excluding those in .gitignore) to the head of the current branch.
// repo: The repo
// message: The commit message
// parentRevs: The revisions corresponding to the commit's parents
func Commit(repo *git.Repository, message string, parentRevs ...string) error {
	// The commit author.
	// TODO: Use an actual user signature
	author := git.Signature{
		"Test User",
		"test@email.com",
		time.Now(),
	}

	// Get the repo's index, which we will use to the stage the files to be committed.
	index, err := repo.Index()
	if err != nil {
		return err
	}

	// Stage all the files in the repo directory (excluding those in .gitignore) for the commit.
	err = index.AddAll(pathSpecs(repo), git.IndexAddDisablePathspecMatch, nil)
	if err != nil {
		return err
	}

	// Write the files in the index into a tree that can be attached to the commit.
	oid, err := index.WriteTree()
	if err != nil {
		return err
	}
	tree, err := repo.LookupTree(oid)
	if err != nil {
		return err
	}

	// Save the index to disk so that it stays in sync with the contents of the working directory.
	// If we don't do this removals of every file are left staged.
	err = index.Write()
	if err != nil {
		return err
	}

	// Retrieve the commit objects associated with the given parent revisions.
	var parentCommits []*git.Commit
	for _, parentRev := range parentRevs {
		parentCommit, err := getCommit(parentRev, repo)
		if err != nil {
			return err
		}

		parentCommits = append(parentCommits, parentCommit)
	}

	// Commit the files to the head of the current branch.
	_, err = repo.CreateCommit("HEAD", &author, &author, message, tree, parentCommits...)
	if err != nil {
		return err
	}

	return nil
}

// Gets the commit corresponding to the given revision
// revision - Revision of the commit to find
// repo - Repo to find the commit in
//
// Returns the commit
func getCommit(revision string, repo *git.Repository) (*git.Commit, error) {
	obj, err := repo.RevparseSingle(revision)
	if err != nil {
		return nil, err
	}
	commit, err := obj.AsCommit()
	if err != nil {
		return nil, err
	}
	return commit, nil
}

// Reverts the last commit WITHOUT leaving a trace of the reverted commit
// reset - If true, the repo is reset back to the last commit
//		   Otherwise, the commit is reverted without resetting the data
func RevertLast(repo *git.Repository, reset bool) error {
	// Gets head commit
	commit, err := getCommit("HEAD", repo)
	if err != nil {return err}

	// Gets commit before head
	oldCommit := commit.Parent(0)
	if oldCommit == nil {return errors.New("head has no parent")}

	// Resets to the last commit
	checkoutOps := git.CheckoutOpts{}
	checkoutOps.Strategy = git.CheckoutForce
	err = repo.ResetToCommit(oldCommit, git.ResetSoft, &checkoutOps)
	if err != nil {return err}

	if reset {
		// Reverts file structure
		index, err := repo.RevertCommit(commit, oldCommit, 0, nil)
		if err != nil {return err}
		err = repo.CheckoutIndex(index, &checkoutOps)
		if err != nil {return err}
	}

	return err
}


func GetLastCommit(repo *git.Repository)  (*git.Commit, error) {
	return getCommit("HEAD", repo)
}