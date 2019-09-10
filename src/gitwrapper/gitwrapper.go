package gitwrapper

import (
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

// Returns the path specs for the Ignore files
func pathSpecs(repo *git.Repository) []string {
	// Finds any ignore files
	ignore := make([]string, 0)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir() + "/.gitignore"), "\n")...)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir() + "/.metroignore"), "\n")...)
	return ignore
}