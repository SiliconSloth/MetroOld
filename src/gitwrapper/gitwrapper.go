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

// Returns the path specs for the Ignore files
func pathSpecs(repo *git.Repository) []string {
	// Finds any ignore files
	ignore := make([]string, 0)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir()+"/.gitignore"), "\n")...)
	ignore = append(ignore, strings.Split(helper.GetFileContents(repo.Workdir()+"/.metroignore"), "\n")...)
	return ignore
}

func SetCreds(repo *git.Repository, username string, email string) error {
	config, err := repo.Config()
	if err != nil { return err }
	err = config.SetString("user.name", username)
	if err != nil { return err }
	err = config.SetString("user.email", email)
	if err != nil { return err }
	return nil
}