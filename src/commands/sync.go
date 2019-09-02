package commands

import "fmt"

func execSync(positionals []string, options map[string]string) {
	fmt.Println(options["timeout"])
}

func printSyncHelp(_ []string, _ map[string]string) {
	fmt.Printf("Usage: metro sync <up | down | <url>>")
}

var Sync = Command{"sync", "Sync with remote repo or something like that", execSync, printSyncHelp}
