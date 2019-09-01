package commands

import "fmt"

func execSync(positionals []string, options map[string]string, hasHelpFlag bool) {
	if hasHelpFlag {
		printSyncHelp()
	} else {
		fmt.Println(options["timeout"])
	}
}

func printSyncHelp() {
	fmt.Printf("Usage: metro sync <up | down | <url>>")
}

var Sync = Command{"sync", "Sync with remote repo or something like that", execSync}
