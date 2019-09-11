package commands

import (
	"fmt"
	git "github.com/libgit2/git2go"
)

func execSync(repo *git.Repository, positionals []string, options map[string]string) error {
	fmt.Println(options["timeout"])
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
		fmt.Println("Could not find given remote")
	}
	remoteStr := remotes[0]
	remote, err := repo.Remotes.Lookup(remoteStr)
	if err != nil { return err }

	fetchOps := git.FetchOptions{}
	err = remote.Fetch(nil, &fetchOps, "Fetch")
	if err != nil { return err }

	return nil
}

func printSyncHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro sync <up | down | <url>>")
}

var Sync = Command{"sync", "Sync with remote repo or something like that", execSync, printSyncHelp}
