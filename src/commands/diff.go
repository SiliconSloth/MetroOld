package commands

func execDiff(positionals []string, options map[string]string, hasHelpFlag bool) {
	println("Diff is not yet implemented.")
}

var Diff = newCommand("diff", "Test out code diff patch functionality", execDiff)
