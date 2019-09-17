package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
)

func execSync(repo *git.Repository, positionals []string, options map[string]string) error {
	var maxArgs int
	var urlIncluded bool
	var url string

	if len(positionals) == 0 {
		maxArgs = 0
		urlIncluded = false
	} else
	if positionals[0] == "down" || positionals[0] == "up" {
		maxArgs = 2
		if len(positionals) == 2 {
			urlIncluded = true
			url = positionals[1]
		} else {
			urlIncluded = false
		}
	} else {
		maxArgs = 1
		urlIncluded = true
		url = positionals[0]
	}

	if len(positionals) > maxArgs {
		return errors.New("Unexpected argument: " + positionals[maxArgs])
	}

	var remote *git.Remote

	remotes, err := repo.Remotes.List()
	if err != nil { return err }

	if urlIncluded {
		remote, err = gitwrapper.AddRemote(repo, url)
		if err != nil { return err }
	} else if len(remotes) < 1 {
		fmt.Println("What url should be fetched from?")
		_, err = fmt.Scan(&url)
		if err != nil { return err }

		remote, err = gitwrapper.AddRemote(repo, url)
		if err != nil { return err }
	} else {
		remote, err = repo.Remotes.Lookup(remotes[0])
		if err != nil { return err }
	}

	if maxArgs == 2 {
		if positionals[0] == "down" {
			_, ok := options["force"]
			err := downsync(repo, remote, ok)
			if err != nil { return err }
		} else
		if positionals[0] == "up" {
			err = upsync(repo, positionals)
			if err != nil { return err }
		} else {
			return errors.New("Unexpected Error: Expected Positional.")
		}
	} else {
		_, ok := options["force"]
		err := downsync(repo, remote, ok)
		if err != nil { return err }
		err = upsync(repo, positionals)
		if err != nil { return err }
	}

	return nil
}

func downsync(repo *git.Repository, remote *git.Remote, force bool) error {
	branch, err := gitwrapper.CurrentBranchName(repo)
	if err != nil { return err }

	callbacks := gitwrapper.CreateCallbacks()
	fetchOps := git.FetchOptions{RemoteCallbacks: callbacks}
	err = remote.Fetch(nil, &fetchOps, "pull")
	if err != nil { return err }

	analysis, err := gitwrapper.MergeAnalysis("origin/master", repo)
	if force && analysis&git.MergeAnalysisUpToDate == 0 {
		err = gitwrapper.ResetHead(repo)
		if err != nil { return err }
	} else if analysis&git.MergeAnalysisUpToDate == 0 {
		areChanged, err := gitwrapper.IsUnsavedChanges(repo)
		if err != nil {
			return err
		}
		if areChanged {
			fmt.Println("Cannot Sync Down with unsaved changes.")
			return nil
		}
	}

	if analysis&git.MergeAnalysisFastForward != 0 {
		err = gitwrapper.FastForward("origin/" + branch, repo)
		if err != nil { return err }
	} else if analysis&git.MergeAnalysisUpToDate == 0 {
		_, err = gitwrapper.CreateBranch(branch + "-local", repo)
		if err != nil { return err }
		err = gitwrapper.FastForward("origin/" + branch, repo)
		if err != nil { return err }

		var in string
		fmt.Println("Conflict Found:")
		fmt.Printf("[0] Absorb local %s line into remote %s line\n", branch, branch)
		fmt.Printf("[1] Move local changes to new line %s-local\n", branch)
		_, err = fmt.Scan(&in)
		if err != nil { return err }
		switch in {
		case "0":
		case "1":
			fmt.Printf("Successfully moved local changes into %s-local\n", branch)
			return nil
		default:
			fmt.Printf("Invalid choice: Moved local changes into %s-local\n", branch)
			return nil
		}

		conflicts, err := gitwrapper.Merge(branch + "-local", repo)
		if err != nil {
			if err.Error() == "Nothing to absorb" {
				return errors.New("You're already in Sync.")
			} else {
				return err
			}
		}

		err = gitwrapper.DeleteBranch(branch + "-local", repo)
		if err != nil { return err }

		if !conflicts {
			err = gitwrapper.Commit(repo, "Completed Absorb.")
			if err != nil { return err }
			fmt.Println("Successfully absorbed changes")
		} else {
			fmt.Println("Conflicts Found: Fix, Commit and Sync again.")
		}
	} else {
		fmt.Println("You're up to date.")
	}

	return nil
}

func upsync(repo *git.Repository, positionals []string) error {
	return nil
}

func printSyncHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro sync <up | down | <url>>")
}

var Sync = Command{"sync", "Sync with remote repo or something like that", execSync, printSyncHelp}
