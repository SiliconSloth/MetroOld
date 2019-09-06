package gitwrapper

import (
	git "github.com/libgit2/git2go"
	"helper"
	"strings"
	"time"
)

// Initialize an empty git repository in the specified directory.
func Init(directory string) (*git.Repository, error) {
	return git.InitRepository(directory+"/.git", false)
}

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

	// Finds any ignore files
	ignore := make([]string, 0)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir() + "/.gitignore"), "\n")...)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir() + "/.metroignore"), "\n")...)

	// Stage all the files in the repo directory (excluding those in .gitignore) for the commit.
	err = index.AddAll(ignore, git.IndexAddDisablePathspecMatch, nil)
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

// Create a new branch from the current head with the specified name.
// Returns the branch
func CreateBranch(name string, repo *git.Repository) (*git.Branch, error) {
	commit, err := getCommit("HEAD", repo)
	if err != nil {
		return nil, err
	}

	branch, err := repo.CreateBranch(name, commit, false)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

// Checks out the given branch by name
// name - Plain Text branch name (e.g. 'master')
// repo - Repo to checkout from
func CheckoutBranch(name string, repo *git.Repository) error {
	err := checkout(name, repo)
	if err != nil {return err}

	err = moveHead(name, repo)
	return err
}

// Moves the head to the given branch
// Files are not changed
func moveHead(name string, repo *git.Repository) error {
	branch, err := repo.LookupBranch(name, git.BranchLocal)
	if err != nil {return err}

	err = repo.SetHead(branch.Reference.Name())
	return err
}

// Checks out the given branch without moving head
// Doesn't change current branch tag
func checkout(name string, repo *git.Repository) error {
	commit, err := getCommit(name, repo)
	if err != nil {return err}
	tree, err := commit.Tree()
	if err != nil {return err}

	checkoutOps := git.CheckoutOpts{}
	checkoutOps.Strategy = git.CheckoutSafe
	err = repo.CheckoutTree(tree, &checkoutOps)

	return err
}