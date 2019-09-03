package main

import (
	"commands"
	"fmt"
	git "github.com/libgit2/git2go"
	"os"
)

func main() {
	success := handleCommand()
	if !success {
		// If no command was given or the command was invalid, show the general help text.
		printHelp()
	}
}

// Gets the current arguments of the program and runs the specified subcommand
//
// Returns: true if the command was found, false if no valid command was specified.
func handleCommand() bool {
	positionals, options, hasHelpFlag, err := commands.ParseArgs(os.Args, allOptions)
	if err != nil {
		// Display the error, print the help text then exit.
		fmt.Println(err.Error())
		return false
	}

	repo, err := git.OpenRepository(".")

	// If we have a sub-command, get it and find the associated command,
	// sending associated data to be executed
	if len(positionals) > 0 {
		argCmd := positionals[0]
		for _, cmd := range allCommands {
			if cmd.Name == argCmd {
				if hasHelpFlag {
					cmd.Help(positionals[1:], options)
				} else {
					// Pass in all positionals after the sub-command.
					err := cmd.Execute(repo, positionals[1:], options)
					if err != nil {
						fmt.Println(err.Error())
						cmd.Help(positionals[1:], options)
					}
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
