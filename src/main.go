package main

import (
	"commands"
	"fmt"
	"os"
)

func main() {
	success := handleCommand()

	if !success {
		// If no command was given or the command or arguments were invalid, show the general help text.
		printHelp()
	}
}

// Gets the current arguments of the program and finds and runs the correct associated subcommand
//
// Returns: True if successful, False if the command fails in any way
func handleCommand() bool {
	positionals, options, hasHelpFlag, err := commands.ParseArgs(os.Args, allOptions)

	if err != nil {
		// Display the error, print the help text then exit.
		fmt.Println(err.Error())
		return false
	}

	// If we have a sub-command, get it and finds the associated command,
	// sending associated data to be executed
	if len(positionals) > 0 {
		argCmd := positionals[0]
		for _, cmd := range allCommands {
			if cmd.Name == argCmd {
				if hasHelpFlag {
					cmd.Help(positionals[1:], options)
				} else {
					// Pass in all positionals after the sub-command.
					cmd.Execute(positionals[1:], options)
				}
				return true
			}
		}
		fmt.Printf("Invalid command: %s\n", argCmd)
	}

	return false
}

// Prints a generic help command with all commands listed
func printHelp() {
	fmt.Println("Usage: metro <command> <args> [options]")
	for _, cmd := range allCommands {
		fmt.Printf("%s - %s\n", cmd.Name, cmd.Description)
	}
	fmt.Println("Use --help for help.")
}