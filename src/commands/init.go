package commands

import (
	"fmt"
	"github.com/libgit2/git2go"
)

func execInit(positionals []string, options map[string]string) {
	println("Initing!")
	_, err := git.InitRepository("test/.git", false) // TODO change to just .git
	if err != nil {
		println(err.Error())
	}
}

func printInitHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro init")
}

var Init = Command{"init", "Test git init", execInit, printInitHelp}
