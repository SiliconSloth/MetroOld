package commands

import (
	"fmt"
	git "github.com/libgit2/git2go"
)

func execTemplate(repo *git.Repository, positionals []string, options map[string]string) error {
	return nil
}

func printTemplateHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro")
}

var Template = Command{"template", "Show you how a command should look", execTemplate, printTemplateHelp}
