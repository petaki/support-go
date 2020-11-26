package cli

import (
	"fmt"
	"strings"
)

type Group struct {
	Name     string
	Usage    string
	Commands []*Command
}

func (g *Group) PrintHelp() int {
	fmt.Println(Yellow("Available commands:"))

	max := 0

	for _, command := range g.Commands {
		if max < len(command.Name) {
			max = len(command.Name)
		}
	}

	max += 2

	for _, command := range g.Commands {
		fmt.Println("  " + Green(command.Name) + strings.Repeat(" ", max-len(command.Name)) + command.Usage)
	}

	return 0
}
