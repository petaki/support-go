package cli

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Command type.
type Command struct {
	Name       string
	Usage      string
	Arguments  []string
	HandleFunc func(group *Group, command *Command, arguments []string) int
	flagSet    *flag.FlagSet
	envErrors  []error
	secrets    map[string]string
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
	err := c.FlagSet().Parse(arguments)
	if err != nil {
		return nil, err
	}

	if len(c.envErrors) > 0 {
		return nil, errors.Join(c.envErrors...)
	}

	err = c.applySecrets()
	if err != nil {
		return nil, err
	}

	if len(c.FlagSet().Args()) < len(c.Arguments) {
		return arguments, ErrMissingArguments
	}

	return c.FlagSet().Args(), nil
}

// PrintHelp function.
func (c *Command) PrintHelp(group *Group) int {
	fmt.Println(Yellow("Usage:"))

	var usage strings.Builder
	usage.WriteString("  " + group.Name + " " + c.Name)

	for _, argument := range c.Arguments {
		usage.WriteString(" <" + argument + ">")
	}

	fmt.Println(usage.String())
	fmt.Println()
	fmt.Println(Yellow("Available flags:"))
	c.FlagSet().PrintDefaults()

	return Success
}

// PrintError function.
func (c *Command) PrintError(err error) int {
	ErrorLog.Output(2, fmt.Sprint(err))

	return Failure
}

// PrintErrorWithHelp function.
func (c *Command) PrintErrorWithHelp(group *Group, err error) int {
	ErrorLog.Output(2, fmt.Sprint(err))
	c.PrintHelp(group)

	return Failure
}
