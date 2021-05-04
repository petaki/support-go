package cli

import (
	"flag"
	"fmt"
)

// Command type.
type Command struct {
	Name       string
	Usage      string
	Arguments  []string
	HandleFunc func(group *Group, command *Command, arguments []string) int
	flagSet    *flag.FlagSet
}

// FlagSet function.
func (c *Command) FlagSet() *flag.FlagSet {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name, flag.ExitOnError)
	}

	return c.flagSet
}

// Parse function.
func (c *Command) Parse(arguments []string) ([]string, error) {
	c.FlagSet().Parse(arguments)

	if len(c.FlagSet().Args()) != len(c.Arguments) {
		return arguments, ErrMissingArguments
	}

	return c.FlagSet().Args(), nil
}

// PrintHelp function.
func (c *Command) PrintHelp(group *Group) int {
	fmt.Println(Yellow("Usage:"))

	usage := "  " + group.Name + " " + c.Name

	for _, argument := range c.Arguments {
		usage += " <" + argument + ">"
	}

	fmt.Println(usage)
	fmt.Println()
	fmt.Println(Yellow("Available flags:"))
	c.FlagSet().PrintDefaults()

	return Success
}

// PrintError function.
func (c *Command) PrintError(err error) int {
	fmt.Println(Red("ERROR\t") + err.Error())

	return Failure
}
