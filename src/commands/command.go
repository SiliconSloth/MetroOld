package commands

import (
	"errors"
	git "github.com/libgit2/git2go"
	"strings"
)

// An Option on the command line.
// Name: The long form Name of the Option, e.g. --help
// Contraction: The short Name of the Option, e.g. -h
// NeedsValue: Whether this Option needs and allows a value associated with it
type Option struct {
	Name        string
	Contraction string
	NeedsValue  bool
}

type Command struct {
	Name        string
	Description string
	Execute     func(*git.Repository, []string, map[string]string)
	Help     	func([]string, map[string]string)
}

// Parse the arguments given to Metro on the command line.
// args: The arguments to parse.
// allOptions: The list of possible options.
//
// All positional arguments must come before all Option arguments.
// Options may have a value associated with them, in the form "--key=value" or "--key value".
// Each Option has a long version, prefixed with --, and a short version, prefixed with -.
// Using the wrong prefix will result in the Option not being recognised.
// The --help and -h flags are excluded from the options; instead hashHelpFlag is set.
// The Name of the executable is excluded from the returned arguments.
//
// Returns: positionals, options, hasHelpFlag, error
// options maps the long Name of each Option to its value, even if a Contraction is used. Prefix -'s are excluded.
//
// Will return an error if:
// - An unknown Option is given
// - An Option that requires a value isn't given one
// - An Option that doesn't require a value is given one
// - A positional argument is found after an Option (interpreted as a value with no corresponding Option flag)
func ParseArgs(args []string, allOptions []Option) ([]string, map[string]string, bool, error) {
	var positionals []string
	options := map[string]string{}
	hasHelpFlag := false

	// Whether we have yet to reach the start of the Option arguments.
	acceptingPositionals := true
	// If the last Option parsed had no value but needed one,
	// this is set so that the next argument can be assigned as its value.
	var openOption Option
	// Whether a Contraction was used to specify the previous Option.
	// Used for printing error messages that match the form used by the user.
	usedContraction := false
	// Start at 1 to exclude the executable Name.
	for _, arg := range args[1:] {
		// If this is an Option flag... (long or contracted, they both start with -)
		if strings.HasPrefix(arg, "-") {
			// Once an Option is found, stop allowing positionals.
			acceptingPositionals = false
			// If the last Option still needs a value, the argument directly after it can't also be an Option Name.
			if openOption.Name != "" {
				return nil, nil, false, optionError("Option needs value", openOption, usedContraction)
			}

			// If this Option had a value specified with --key=value, extract it.
			// If not the value will be "" and should be assigned by the next argument.
			key, value := splitAtFirst(arg, "=")

			var opt Option
			usedContraction = !strings.HasPrefix(key, "--")
			// Find which Option this key corresponds to.
			if usedContraction {
				key = key[1:] // Remove -
				for _, o := range allOptions {
					if o.Contraction == key {
						opt = o
						break
					}
				}
			} else {
				key = key[2:] //Remove --
				for _, o := range allOptions {
					if o.Name == key {
						opt = o
						break
					}
				}
			}

			if opt.Name == "" {
				return nil, nil, false, optionError("Bad Option", Option{key, key, false}, usedContraction)
			}

			if value != "" {
				if opt.NeedsValue {
					options[opt.Name] = value
				} else {
					return nil, nil, false, optionError("Option doesn't take a value", opt, usedContraction)
				}
			} else if opt.NeedsValue {
				// If no value was given but this Option needs one, consume the next argument as the value.
				openOption = opt
			} else {
				if opt.Name == "help" {
					// Rather than including the --help flag in the options, set hasHelpFlag.
					hasHelpFlag = true
				} else {
					// Option is present but has no value.
					options[opt.Name] = ""
				}
			}
		} else if openOption.Name != "" {
			//  If the last Option had no value provided with = but needed one, this argument becomes its value.
			options[openOption.Name] = arg
			// Reset to an empty Option.
			openOption = Option{}
		} else {
			if acceptingPositionals {
				positionals = append(positionals, arg)
			} else {
				// If a positional argument is found after the options have started,
				// we assume it is a value missing an Option key.
				return nil, nil, false, errors.New("Value without flag: " + arg)
			}
		}
	}

	// Make sure the last Option had a value.
	if openOption.Name != "" {
		return nil, nil, false, optionError("Option needs a value", openOption, usedContraction)
	}
	return positionals, options, hasHelpFlag, nil
}

// Split the given string around the fist occurrence of the given substring.
// If the substring is not found, the input string is returned as the first string and "" as the second.
func splitAtFirst(str string, sub string) (string, string) {
	index := strings.Index(str, sub)
	if index == -1 {
		return str, ""
	} else {
		return str[:index], str[index+len(sub):]
	}
}

// Create a properly formatted error message about an Option being inputted incorrectly,
// taking into account whether the user used the contracted form of the Name or not.
func optionError(message string, opt Option, usedContraction bool) error {
	if usedContraction {
		return errors.New(message + ": -" + opt.Contraction)
	} else {
		return errors.New(message + ": --" + opt.Name)
	}
}
