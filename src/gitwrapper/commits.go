package gitwrapper

import (
	"errors"
	git "github.com/libgit2/git2go"
	"helper"
	"strings"
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
// reset - If true, commit is deleted and working directory reset to last commit
//		   Otherwise working directory is unchanged
func RevertLastCommit(repo *git.Repository, reset bool) error {
	return RevertCommit(repo, 1, reset)
}

// Reverts the last commit WITHOUT leaving a trace of the reverted commit
// commitsBack - How many commits back to revert
// reset - If true, commit is deleted and working directory reset to last commit
//		   Otherwise working directory is unchanged
func RevertCommit(repo *git.Repository, commitsBack int, reset bool) error {
	if commitsBack < 1 {
		return errors.New("Invalid commit to delete.")
	}

	// Gets head commit
	commit, err := getCommit("HEAD", repo)
	if err != nil {
		return err
	}

	// Gets commit before head
	oldCommit := commit
	for ; commitsBack > 0; commitsBack-- {
		oldCommit = oldCommit.Parent(0)
		if oldCommit == nil {
			return errors.New("head has no parent")
		}
	}

	// Resets head to the last commit, deleting the current head
	// If reset is true, also resets working directory
	checkoutOps := git.CheckoutOpts{}
	checkoutOps.Strategy = git.CheckoutForce
	var resetType git.ResetType
	if reset {
		resetType = git.ResetHard
	} else {
		resetType = git.ResetSoft
	}
	err = repo.ResetToCommit(oldCommit, resetType, &checkoutOps)
	if err != nil {
		return err
	}

	return err
}

func GetLastCommit(repo *git.Repository) (*git.Commit, error) {
	return getCommit("HEAD", repo)
}

// If anything is added, creates a new branch with a commit called WIP
func WIPCommit(repo *git.Repository) error {
	statusOps := git.StatusOptions{
		Show:     git.StatusShowIndexAndWorkdir,
		Flags:    git.StatusOptIncludeUntracked,
		Pathspec: pathSpecs(repo),
	}
	status, err := repo.StatusList(&statusOps)
	if err != nil {
		return err
	}
	count, err := status.EntryCount()
	if err != nil {
		return err
	}

	// If nothing to commit, don't bother with a WIP
	if count == 0 {
		return nil
	}

	name, err := CurrentBranchName(repo)
	if err != nil {
		return err
	}
	if strings.HasSuffix(name, helper.WipString) {
		return nil
	}

	// If WIP already exists, delete
	if BranchExists(name+helper.WipString, repo) {
		err = DeleteBranch(name+helper.WipString, repo)
		if err != nil {
			return err
		}
	}

	_, err = CreateBranch(name+helper.WipString, repo)
	if err != nil {
		return err
	}
	err = moveHead(name+helper.WipString, repo)
	if err != nil {
		return err
	}

	merging := MergeOngoing(repo)
	if merging {
		err = Commit(repo, "WIP", "HEAD^{commit}", "MERGE_HEAD^{commit}")
		if err != nil {
			return err
		}
		err = repo.StateCleanup()
	} else {
		err = Commit(repo, "WIP", "HEAD^{commit}")
	}
	if err != nil {
		return err
	}

	return nil
}

// Deletes the WIP commit at head if any
func WIPUncommit(repo *git.Repository) error {
	name, err := CurrentBranchName(repo)
	if err != nil {
		return err
	}

	// No WIP branch
	if !BranchExists(name+helper.WipString, repo) {
		return nil
	}

	commit, err := getCommit(name+helper.WipString, repo)
	if err != nil {
		return err
	}
	if commit.ParentCount() > 1 {
		mergeHead := commit.Parent(1).Id().String()
		conflicts, err := Merge(mergeHead, repo)
		if err != nil {
			return err
		}
		if !conflicts {
			return errors.New("WIP contained merge with no conflicts.")
		}
	}

	err = checkout(name+helper.WipString, false, repo)
	if err != nil {
		return err
	}

	err = DeleteBranch(name+helper.WipString, repo)
	if err != nil {
		return err
	}

	return nil
}
