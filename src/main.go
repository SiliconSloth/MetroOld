package main

import (
	"commands"
	"fmt"
	"os"
)

func main() {
	positionals, options, hasHelpFlag, err := commands.ParseArgs(os.Args, allOptions)
	if err != nil {
		// Display the error, print the help text then exit.
		fmt.Println(err.Error())
		// If we actually have a sub-command...
	} else if len(positionals) > 0 {
		argCmd := positionals[0]
		for _, cmd := range allCommands {
			if cmd.Name == argCmd {
				cmd.Execute(positionals[1:], options, hasHelpFlag) // Pass in all positionals after the sub-command.
				return                                             // Don't print the help text below.
			}
		}
		fmt.Printf("Invalid command: %s\n", argCmd)
	}
	// If no command was given or the command or arguments were invalid, show the general help text.
	printHelp()
}

func printHelp() {
	fmt.Println("Usage: metro <command> <args> [options]")
	for _, cmd := range allCommands {
		fmt.Printf("%s - %s\n", cmd.Name, cmd.Description)
	}
	fmt.Println("Use --help for help.")
}
