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
			err := downsync(repo, remote)
			if err != nil { return err }
		} else
		if positionals[0] == "up" {
			err = upsync(repo, positionals)
			if err != nil { return err }
		} else {
			return errors.New("Unexpected Error: Expected Positional.")
		}
	} else {
		err := downsync(repo, remote)
		if err != nil { return err }
		err = upsync(repo, positionals)
		if err != nil { return err }
	}

	return nil
}

func downsync(repo *git.Repository, remote *git.Remote) error {
	callbacks := gitwrapper.CreateCallbacks()
	fetchOps := git.FetchOptions{ RemoteCallbacks: callbacks }
	err := remote.Fetch(nil, &fetchOps, "pull")
	if err != nil { return err }

	branch, err := gitwrapper.CurrentBranchName(repo)
	if err != nil { return err }
	conflicts, err := gitwrapper.Merge("origin/" + branch, repo)
	if err != nil {
		if err.Error() == "Nothing to absorb" {
			return errors.New("You're already in Sync.")
		} else { return err }
	}

	if !conflicts {
		fmt.Println("Successfully Downsynched.")
	} else {
		fmt.Println("Conflicts Found: Fix, Commit and Sync again.")
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
