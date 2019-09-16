package gitwrapper

import (
	"fmt"
	git "github.com/libgit2/git2go"
)

// Adds the origin remote if it doesn't exist
// Edits the url of the existing origin remote if it doesn't
func AddRemote(repo *git.Repository, url string) (*git.Remote, error){
	var remote *git.Remote

	remotes, err := repo.Remotes.List()
	if err != nil { return nil, err }

	if len(remotes) < 1 {
		remote, err = repo.Remotes.Create("origin", url)
		if err != nil { return nil, err }

		return remote, nil
	} else {
		err = repo.Remotes.SetUrl("origin", url)
		if err != nil { return nil, err }

		return repo.Remotes.Lookup(remotes[0])
	}
}

// Creates the variable holding the Callbacks
func CreateCallbacks() git.RemoteCallbacks {
	callbacks := git.RemoteCallbacks{}
	callbacks.TransferProgressCallback = transferProgressCallback
	callbacks.CredentialsCallback = credentialsCallback
	return callbacks
}

func transferProgressCallback(stats git.TransferProgress) git.ErrorCode {
	progress := (100 * (stats.ReceivedObjects + stats.IndexedObjects)) / (2 * stats.TotalObjects)
	fmt.Printf("\rProgress: %d%%", progress)
	if progress == 100 { fmt.Println() }
	return git.ErrOk
}

func credentialsCallback(url string, username_from_url string, allowed_types git.CredType) (code git.ErrorCode, cred *git.Cred) {
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