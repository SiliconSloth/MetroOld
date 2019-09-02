package commands

import (
	"fmt"
	git "github.com/libgit2/git2go"
)

func execDiff(_ *git.Repository, positionals []string, options map[string]string) {
	println("Diff is not yet implemented.")
}

func printDiffHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro diff <file1> <file2>")
}

var Diff = Command{"diff", "Test out code diff patch functionality", execDiff, printDiffHelp}
