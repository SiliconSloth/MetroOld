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

func AssertConflicts(repo *git.Repository) error {
	index, err := repo.Index()
	if err != nil {
		return err
	}

	if index.HasConflicts() {
		return errors.New("Branch has conflicts, please finish reosolving them.")
	}
	return nil
}

// Returns the path specs for the Ignore files
func pathSpecs(repo *git.Repository) []string {
	// Finds any ignore files
	ignore := make([]string, 0)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir()+"/.gitignore"), "\n")...)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir()+"/.metroignore"), "\n")...)
	return ignore
}
