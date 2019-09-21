package metro

import (
	"errors"
	git "github.com/libgit2/git2go"
	"io/ioutil"
	"strings"
)

func Absorb(mergeHead string, repo *git.Repository) (bool, error) {
	if strings.HasSuffix(mergeHead, WipString) {
		return false, errors.New("Can't absorb WIP branch.")
	}

	err := AssertMerging(repo)
	if err != nil {
		return false, err
	}

	err = startMerge(mergeHead, repo)
	if err != nil {
		return false, err
	}

	index, err := repo.Index()
	if err != nil {
		return false, err
	}

	if index.HasConflicts() {
		return true, nil
	} else {
		// If no conflicts occurred make the merge commit right away.
		err = Resolve(repo)
		if err != nil {
			return false, err
		}
		return false, nil
	}
}

// Merge the specified commit into the current branch head.
// The repo will be left in a merging state, possibly with conflicts in the index.
func startMerge(name string, repo *git.Repository) error {
	otherHead, err := GetCommit(name, repo)
	if err != nil {
		return err
	}
	annOther, err := repo.LookupAnnotatedCommit(otherHead.Id())
	if err != nil {
		return err
	}
	sources := []*git.AnnotatedCommit{annOther}

	analysis, _, err := repo.MergeAnalysis(sources)
	if err != nil {
		return err
	}
	if analysis&git.MergeAnalysisNone != 0 || analysis&git.MergeAnalysisUpToDate != 0 {
		return errors.New("Nothing to absorb.")
	}
	if analysis&git.MergeAnalysisNormal == 0 {
		return errors.New("Non-normal absorb.")
	}

	mergeOptions, err := git.DefaultMergeOptions()
	if err != nil {
		return err
	}
	checkoutOptions := git.CheckoutOpts{
		Strategy: git.CheckoutForce | git.CheckoutAllowConflicts,
	}

	err = repo.Merge(sources, &mergeOptions, &checkoutOptions)
	if err != nil {
		return err
	}
	err = setMergeMessage(defaultMergeMessage(name), repo)
	if err != nil {
		return err
	}

	return nil
}

// Create a commit of the ongoing merge and clear the merge state and conflicts from the repo.
func Resolve(repo *git.Repository) error {
	merging := MergeOngoing(repo)
	if !merging {
		return errors.New("You can only resolve conflicts while absorbing.")
	}

	mergedID, err := mergeHeadID(repo)
	if err != nil {
		return err
	}

	message, err := getMergeMessage(repo)
	if err != nil {
		return err
	}

	// Remove merge state.
	err = repo.StateCleanup()
	if err != nil {
		return err
	}

	// Remove index conflicts.
	index, err := repo.Index()
	if err != nil {
		return err
	}
	index.CleanupConflicts()

	err = Commit(repo, message, "HEAD^{commit}", mergedID+"^{commit}")
	if err != nil {
		return err
	}

	return nil
}

// Get the commit ID of the merge head. Assumes a merge is ongoing.
func mergeHeadID(repo *git.Repository) (string, error) {
	mergeHead, err := GetCommit("MERGE_HEAD^{commit}", repo)
	if err != nil {
		return "", err
	}
	return mergeHead.Id().String(), nil
}

// The commit message Metro uses when absorbing a commit referenced by the given name.
func defaultMergeMessage(mergedName string) string {
	return "Absorbed " + mergedName
}

// Load the merge message from the repo's MERGE_MSG file.
func getMergeMessage(repo *git.Repository) (string, error) {
	dat, err := ioutil.ReadFile(repo.Path() + "/MERGE_MSG")
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// Write a merge message to the repo's MERGE_MSG file.
func setMergeMessage(message string, repo *git.Repository) error {
	return ioutil.WriteFile(repo.Path()+"/MERGE_MSG", []byte(message), 0644)
}
