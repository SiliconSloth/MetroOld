package commands

import (
	"fmt"
	git "github.com/libgit2/git2go"
)

func execTemplate(_ *git.Repository, positionals []string, options map[string]string) {

}

func printTemplateHelp(_ []string, _ map[string]string) {
	fmt.Println("Usage: metro")
}

var Template = Command{"template", "Template Command", execTemplate, printTemplateHelp}

