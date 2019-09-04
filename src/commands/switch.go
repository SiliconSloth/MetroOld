package commands

import (
"fmt"
git "github.com/libgit2/git2go"
)

func execSwitch(_ *git.Repository, positionals []string, options map[string]string) error {
	return nil
}

func printSwitchHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro switch <line>")
}

var Switch = Command{"switch", "Switch Command", execSwitch, printSwitchHelp}
