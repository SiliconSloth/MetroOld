package commands

import (
	"fmt"
	git "github.com/libgit2/git2go"
)

func execSync(_ *git.Repository, positionals []string, options map[string]string) error {
	fmt.Println(options["timeout"])
	return nil
}

func printSyncHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro sync <up | down | <url>>")
}

var Sync = Command{"sync", "Sync with remote repo or something like that", execSync, printSyncHelp}
