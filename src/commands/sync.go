package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
)

func execSync(repo *git.Repository, positionals []string, options map[string]string) error {
	//_, cred := git.NewCredDefault()
	remotes, err := repo.Remotes.List()
	if err != nil { return err }
	if len(remotes) < 1 {
		fmt.Println("What url should be fetched from?")
		var url string
		_, err = fmt.Scan(&url)
		if err != nil { return err }

		err = repo.Remotes.AddFetch(url, "origin")
		if err != nil { return err }

		remotes, err = repo.Remotes.List()
	}
	if len(remotes) < 1 {
		return errors.New("Could not find given remote.")
	}
	remoteStr := remotes[0]
	err = repo.Remotes.SetUrl("origin", "https://github.com/Black-Photon/Metro-Test")
	if err != nil { return err }
	remote, err := repo.Remotes.Lookup(remoteStr)
	if err != nil { return err }

	callbacks := git.RemoteCallbacks{}
	callbacks.TransferProgressCallback = func(stats git.TransferProgress) git.ErrorCode {
		println("Indexed Objects:", stats.IndexedObjects)
		println("Local Objects:", stats.LocalObjects)
		println("Received Bytes:", stats.ReceivedBytes)
		println("Received Objects:", stats.ReceivedObjects)
		println("Total Deltas:", stats.TotalDeltas)
		println("Total Objects:", stats.TotalObjects)

		return git.ErrOk
	}
	callbacks.CredentialsCallback = func(url string, username_from_url string, allowed_types git.CredType) (code git.ErrorCode, cred *git.Cred) {
		_, creds := git.NewCredUserpassPlaintext("Black-Photon", "")
		return git.ErrOk, &creds
	}
	fetchOps := git.FetchOptions{ RemoteCallbacks: callbacks }
	err = remote.Fetch(nil, &fetchOps, "pull")
	if err != nil { return err }

	conflicts, err := gitwrapper.Merge("origin/master", repo)
	if err != nil { return err }

	if !conflicts {
		println("Successfully Downsynched")
	} else {
		println("Conflicts Found")
	}

	return nil
}

func printSyncHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro sync <up | down | <url>>")
}

var Sync = Command{"sync", "Sync with remote repo or something like that", execSync, printSyncHelp}
