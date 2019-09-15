package commands

import (
	"errors"
	"fmt"
	git "github.com/libgit2/git2go"
	"gitwrapper"
)

func execSync(repo *git.Repository, positionals []string, options map[string]string) error {
	if len(positionals) > 1 {
		return errors.New("Unexpected argument: " + positionals[1])
	}

	var remote *git.Remote

	remotes, err := repo.Remotes.List()
	if err != nil { return err }

	if len(positionals) == 1 {
		if len(remotes) < 1 {
			remote, err = repo.Remotes.Create("origin", positionals[0])
		} else {
			err = repo.Remotes.SetUrl("origin", positionals[0])
		}
		if err != nil {
			return err
		}
	} else if len(remotes) < 1 {
		fmt.Println("What url should be fetched from?")
		var url string
		_, err = fmt.Scan(&url)
		if err != nil { return err }

		_, err = repo.Remotes.Create("origin", url)
		if err != nil { return err }

		remotes, err = repo.Remotes.List()

		if len(remotes) < 1 {
			return errors.New("Could not find given remote.")
		}

		remoteStr := remotes[0]
		remote, err = repo.Remotes.Lookup(remoteStr)
		if err != nil { return err }
	}

	callbacks := git.RemoteCallbacks{}
	callbacks.TransferProgressCallback = func(stats git.TransferProgress) git.ErrorCode {
		progress := (100 * (stats.ReceivedObjects + stats.IndexedObjects)) / (2 * stats.TotalObjects)
		fmt.Printf("\rProgress: %d%%", progress)
		if progress == 100 { fmt.Println() }
		return git.ErrOk
	}
	callbacks.CredentialsCallback = func(url string, username_from_url string, allowed_types git.CredType) (code git.ErrorCode, cred *git.Cred) {
		var err int
		var creds git.Cred
		switch allowed_types {
		case git.CredTypeDefault:
			err, creds = git.NewCredDefault()
		case git.CredTypeUserpassPlaintext:
			fmt.Print("Username: ")
			var username string
			_, err1 := fmt.Scan(&username)
			if err1 != nil {
				fmt.Println(err1.Error())
				return git.ErrGeneric, nil
			}

			fmt.Print("Password: ")
			var password string
			_, err1 = fmt.Scan(&password)
			if err1 != nil {
				fmt.Println(err1.Error())
				return git.ErrGeneric, nil
			}

			err, creds = git.NewCredUserpassPlaintext(username, password)
		default:
			fmt.Println("Metro currently doesn't support SSH. Please use HTML.")
			return git.ErrGeneric, nil
		}
		if err != 0 {
			return git.ErrGeneric, nil
		}
		return git.ErrOk, &creds
	}
	fetchOps := git.FetchOptions{ RemoteCallbacks: callbacks }
	err = remote.Fetch(nil, &fetchOps, "pull")
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

func printSyncHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro sync <up | down | <url>>")
}

var Sync = Command{"sync", "Sync with remote repo or something like that", execSync, printSyncHelp}
