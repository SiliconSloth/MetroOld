package metro

import (
	"errors"
	git "github.com/libgit2/git2go"
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

	err = index.AddAll(nil, git.IndexAddDisablePathspecMatch, nil)
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
		parentCommit, err := GetCommit(parentRev, repo)
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
func GetCommit(revision string, repo *git.Repository) (*git.Commit, error) {
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

// TODO: Make this work with commits with more than one parent
func Patch(repo *git.Repository, message string) error {
	err := AssertMerging(repo)
	if err != nil {
		return err
	}

	err = DeleteLastCommit(repo, false)
	if err != nil {
		return err
	}

	err = Commit(repo, message, "HEAD^{commit}")
	if err != nil {
		return err
	}
	return nil
}

// Reverts the last commit WITHOUT leaving a trace of the reverted commit
// reset - If true, commit is deleted and working directory reset to last commit
//		   Otherwise working directory is unchanged
func DeleteLastCommit(repo *git.Repository, reset bool) error {
	return DeleteCommits(repo, 1, reset)
}

// Reverts the last commit WITHOUT leaving a trace of the reverted commit
// commitsBack - How many commits back to revert
// reset - If true, commit is deleted and working directory reset to last commit
//		   Otherwise working directory is unchanged
func DeleteCommits(repo *git.Repository, commitsBack int, reset bool) error {
	if commitsBack < 1 {
		return errors.New("Invalid commit to delete.")
	}

	// Gets head commit
	commit, err := GetCommit("HEAD", repo)
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

// Checks out the given commit without moving head,
// such that the working directory will match the commit contents.
// Doesn't change current branch ref.
func checkout(name string, repo *git.Repository) error {
	commit, err := GetCommit(name, repo)
	if err != nil {
		return err
	}
	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	checkoutOps := git.CheckoutOpts{}
	checkoutOps.Strategy = git.CheckoutForce
	err = repo.CheckoutTree(tree, &checkoutOps)
	if err != nil {
		return err
	}

	return err
}

func CommitExists(name string, repo *git.Repository) bool {
	_, err := GetCommit(name, repo)
	return err == nil
}

// If anything is added, creates a new branch with a commit called WIP
func SaveWIP(repo *git.Repository) error {
	statusOps := git.StatusOptions{
		Show:  git.StatusShowIndexAndWorkdir,
		Flags: git.StatusOptIncludeUntracked,
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
	if strings.HasSuffix(name, WipString) {
		return nil
	}

	// If WIP already exists, delete
	if CommitExists(name+WipString, repo) {
		err = DeleteBranch(name+WipString, repo)
		if err != nil {
			return err
		}
	}

	_, err = CreateBranch(name+WipString, repo)
	if err != nil {
		return err
	}
	err = moveHead(name+WipString, repo)
	if err != nil {
		return err
	}

	merging := MergeOngoing(repo)
	if merging {
		message, err := getMergeMessage(repo)
		if err != nil {
			return err
		}

		// Store the merge message in the second line (and beyond) of the WIP commit message.
		err = Commit(repo, "WIP\n"+message, "HEAD^{commit}", "MERGE_HEAD^{commit}")
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

// Deletes the WIP commit at head if any, restoring the contents to the working directory
// and resuming a merge if one was ongoing.
func RestoreWIP(repo *git.Repository) error {
	name, err := CurrentBranchName(repo)
	if err != nil {
		return err
	}

	// No WIP branch
	if !CommitExists(name+WipString, repo) {
		return nil
	}

	wipCommit, err := GetCommit(name+WipString, repo)
	if err != nil {
		return err
	}

	index, err := repo.Index()
	if err != nil {
		return err
	}

	var conflicts []git.IndexConflict
	// If the WIP commit has two parents a merge was ongoing.
	if wipCommit.ParentCount() > 1 {
		mergeHead := wipCommit.Parent(1).Id().String()
		err := startMerge(mergeHead, repo)
		if err != nil {
			return err
		}

		// Reload the merge message from before, stored in the second line (and beyond)
		// of the WIP commit message.
		commitMessage := wipCommit.Message()
		newlineIndex := strings.Index(commitMessage, "\n")
		// If the commit message only has one line (only happens if it has been tampered with)
		// just leave the message as the default one created when restarting the merge.
		// Otherwise restore the merge message from the commit message.
		if newlineIndex >= 0 {
			mergeMessage := commitMessage[newlineIndex+1:]
			err = setMergeMessage(mergeMessage, repo)
			if err != nil {
				return err
			}
		}

		// Remove the conflicts from the index temporarily so we can checkout.
		// They will be restored after so that the index and working dir
		// match their state when the WIP commit was created.
		conflicts, err = getConflicts(index)
		if err != nil {
			return err
		}
		index.CleanupConflicts()
	}

	// Restore the contents of the WIP commit to the working directory.
	err = checkout(name+WipString, repo)
	if err != nil {
		return err
	}

	err = DeleteBranch(name+WipString, repo)
	if err != nil {
		return err
	}

	// If we are mid-merge, restore the conflicts from the merge.
	for _, conflict := range conflicts {
		err = index.AddConflict(conflict.Ancestor, conflict.Our, conflict.Their)
		if err != nil {
			return err
		}
	}
	err = index.Write()
	if err != nil {
		return err
	}

	return nil
}
