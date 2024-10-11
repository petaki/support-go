package cli

import (
	"fmt"
	"strings"
)

// Group type.
type Group struct {
	Name     string
	Usage    string
	Commands []*Command
}

// PrintHelp function.
func (g *Group) PrintHelp() int {
	fmt.Println(Yellow("Available commands:"))

	maxLength := 0

	for _, command := range g.Commands {
		if maxLength < len(command.Name) {
			maxLength = len(command.Name)
		}
	}

	maxLength += 2

	for _, command := range g.Commands {
		fmt.Println("  " + Green(command.Name) + strings.Repeat(" ", maxLength-len(command.Name)) + command.Usage)
	}

	return Success
}
