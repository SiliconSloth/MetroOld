package gitwrapper

import (
	"errors"
	git "github.com/libgit2/git2go"
	"helper"
	"strings"
)

// Initialize an empty git repository in the specified directory.
func Init(directory string) (*git.Repository, error) {
	return git.InitRepository(directory+"/.git", false)
}

// If anything is added, creates a new branch with a commit called WIP
func WIPCommit(repo *git.Repository) error {
	statusOps := git.StatusOptions{
		Show: git.StatusShowIndexAndWorkdir,
		Flags: git.StatusOptIncludeUntracked,
		Pathspec: pathSpecs(repo),
	}
	status, err := repo.StatusList(&statusOps)
	if err != nil {return err}
	count, err := status.EntryCount()
	if err != nil {return err}

	// If nothing to commit, don't bother with a WIP
	if count == 0 {
		return nil
	}

	name, err := currentBranchName(repo)
	if err != nil {return err}
	if strings.HasSuffix(name, helper.WipString) {
		return nil
	}

	// If WIP already exists, delete
	if branchExists(name + helper.WipString, repo) {
		err = deleteBranch(name + helper.WipString, repo)
		if err != nil {return err}
	}

	_, err = CreateBranch(name + helper.WipString, repo)
	if err != nil {return err}
	err = moveHead(name + helper.WipString, repo)
	if err != nil {return err}
	err = Commit(repo, "WIP", "HEAD^{commit}")
	if err != nil {return err}

	return nil
}

// Deletes the WIP commit at head if any
func WIPUncommit(repo *git.Repository) error {
	name, err := currentBranchName(repo)
	if err != nil {return err}

	// No WIP branch
	if !branchExists(name + helper.WipString, repo) {
		return nil
	}
	err = checkout(name + helper.WipString, repo)
	if err != nil {return err}

	err = deleteBranch(name + helper.WipString, repo)
	if err != nil {return err}

	return nil
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

// Returns the path specs for the Ignore files
func pathSpecs(repo *git.Repository) []string {
	// Finds any ignore files
	ignore := make([]string, 0)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir() + "/.gitignore"), "\n")...)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir() + "/.metroignore"), "\n")...)
	return ignore
}