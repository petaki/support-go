package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// String function.
func (c *Command) String(name, env, value, usage string) *string {
	if raw := os.Getenv(env); raw != "" {
		value = raw
	}

	return c.FlagSet().String(name, value, envUsage(usage, env))
}

// Int function.
func (c *Command) Int(name, env string, value int, usage string) *int {
	if raw := os.Getenv(env); raw != "" {
		number, err := strconv.Atoi(raw)
		if err != nil {
			c.envErrors = append(c.envErrors, fmt.Errorf("%w: %s, %s, %w", ErrInvalidEnvValue, env, raw, err))
		} else {
			value = number
		}
	}

	return c.FlagSet().Int(name, value, envUsage(usage, env))
}

// Bool function.
func (c *Command) Bool(name, env string, value bool, usage string) *bool {
	if raw := os.Getenv(env); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			c.envErrors = append(c.envErrors, fmt.Errorf("%w: %s, %s, %w", ErrInvalidEnvValue, env, raw, err))
		} else {
			value = parsed
		}
	}

	return c.FlagSet().Bool(name, value, envUsage(usage, env))
}

// Secret function.
func (c *Command) Secret(name, env, usage string) *string {
	if c.secrets == nil {
		c.secrets = make(map[string]string)
	}

	c.secrets[name] = env

	return c.FlagSet().String(name, "", usage+" (default $"+env+")")
}

func (c *Command) applySecrets() error {
	visited := make(map[string]bool)

	c.FlagSet().Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})

	for name, env := range c.secrets {
		value := os.Getenv(env)
		if visited[name] || value == "" {
			continue
		}

		err := c.FlagSet().Set(name, value)
		if err != nil {
			return err
		}
	}

	return nil
}
