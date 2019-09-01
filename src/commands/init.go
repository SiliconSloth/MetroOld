package commands

import "github.com/libgit2/git2go"

func execInit(positionals []string, options map[string]string, hasHelpFlag bool) {
	println("Initing!")
	_, err := git.InitRepository("test/.git", true)
	if err != nil {
		println(err.Error())
	}
}

var Init = Command{"init", "Test git init", execInit}
