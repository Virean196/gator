package main

import (
	"fmt"
	"os"

	"github.com/virean196/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	value, exists := c.handlers[cmd.name]
	if exists {
		err := value(s, cmd)
		if err != nil {
			return fmt.Errorf("error running the function")
		}
		return err
	} else {
		return fmt.Errorf("command doesnt exist")
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) <= 1 {
		fmt.Print("username is required\n")
		os.Exit(1)
	}
	s.cfg.SetUser(cmd.args[1])
	fmt.Printf("User: %s has been set\n", cmd.args[1])
	return nil
}

func main() {
	var s state
	s.cfg, _ = config.Read()
	commands := commands{handlers: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Print("not enough args ")
		os.Exit(1)
	}
	command := command{name: args[1], args: args[1:]}
	commands.run(&s, command)
}
