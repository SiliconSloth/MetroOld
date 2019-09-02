package commands

import "fmt"

func execDiff(positionals []string, options map[string]string) {
	println("Diff is not yet implemented.")
}

func printDiffHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro diff <file1> <file2>")
}

var Diff = Command{"diff", "Test out code diff patch functionality", execDiff, printDiffHelp}
