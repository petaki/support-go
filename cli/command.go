package cli

import (
	"errors"
	"flag"
	"fmt"
)

type Command struct {
	Name       string
	Usage      string
	Arguments  []string
	HandleFunc func(group *Group, command *Command, arguments []string) int
	flagSet    *flag.FlagSet
}

func (c *Command) FlagSet() *flag.FlagSet {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name, flag.ExitOnError)
	}

	return c.flagSet
}

func (c *Command) Parse(arguments []string) ([]string, error) {
	c.FlagSet().Parse(arguments)

	if len(c.FlagSet().Args()) != len(c.Arguments) {
		return arguments, errors.New("missing arguments")
	}

	return c.FlagSet().Args(), nil
}

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

	return 0
}
