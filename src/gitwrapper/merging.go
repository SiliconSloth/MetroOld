package gitwrapper

import (
	"errors"
	git "github.com/libgit2/git2go"
)

// Merge the specified commit into the current branch head.
// The repo will be left in a merging state, possibly with conflicts in the index.
func StartMerge(name string, repo *git.Repository) error {
	otherHead, err := getCommit(name, repo)
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

	return nil
}

// Create a commit of the ongoing merge and clear the merge state and conflicts from the repo.
func MergeCommit(repo *git.Repository) error {
	mergeHead, err := getCommit("MERGE_HEAD^{commit}", repo)
	if err != nil {
		return err
	}
	mergedID := mergeHead.Id().String()

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

	err = Commit(repo, "Absorbed "+mergedID, "HEAD^{commit}", mergedID+"^{commit}")
	if err != nil {
		return err
	}

	return nil
}
