package metro

import (
	"errors"
	git "github.com/libgit2/git2go"
)

const (
	WipString = "#wip"
)

// Initialize an empty git repository in the specified directory.
func Create(directory string) (*git.Repository, error) {
	repo, err := git.InitRepository(directory+"/.git", false)
	if err != nil {
		return nil, err
	}

	err = Commit(repo, "Create repository")
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// Raises an error if the repo is currently in merging state.
func AssertMerging(repo *git.Repository) error {
	if MergeOngoing(repo) {
		return errors.New("Branch has conflicts, please finish resolving them.\nRun metro resolve when you are done.")
	}
	return nil
}

// Returns true if the repo is currently in merging state.
func MergeOngoing(repo *git.Repository) bool {
	_, err := repo.RevparseSingle("MERGE_HEAD")
	return err == nil
}

// Return a list of all the conflicts in the given index.
func getConflicts(index *git.Index) ([]git.IndexConflict, error) {
	iterator, err := index.ConflictIterator()
	if err != nil {
		return nil, err
	}

	var conflicts []git.IndexConflict
	for {
		conflict, err := iterator.Next()
		if git.IsErrorCode(err, git.ErrIterOver) {
			break
		}
		if err != nil {
			return nil, err
		}
		conflicts = append(conflicts, conflict)
	}
	return conflicts, nil
}
