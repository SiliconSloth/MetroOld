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

func AssertMerging(repo *git.Repository) error {
	if MergeOngoing(repo) {
		return errors.New("Branch has conflicts, please finish reosolving them.")
	}
	return nil
}

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

// Returns the path specs for the Ignore files
func pathSpecs(repo *git.Repository) []string {
	// Finds any ignore files
	ignore := make([]string, 0)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir()+"/.gitignore"), "\n")...)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir()+"/.metroignore"), "\n")...)
	return ignore
}
