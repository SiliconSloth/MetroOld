package commands

import (
	"fmt"
	git "github.com/libgit2/git2go"
)

func execSwitch(_ *git.Repository, _ []string, _ map[string]string) error {
	return nil
}

func printSwitchHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro switch <line>")
}

var Switch = Command{"switch", "Switch to a different line or branch", execSwitch, printSwitchHelp}
