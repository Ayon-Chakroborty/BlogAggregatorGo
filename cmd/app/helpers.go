package main

import (
	"encoding/json"
	"os"
)

type GatorCommand func(Command) error

type Command struct {
	Name string
	Args []string
}

func (a *Application) Run() error {
	args := os.Args[1:]
	fn, err := a.CheckGatorCommand(args)
	if err != nil {
		return err
	}

	cmd := Command{
		Name: args[1],
		Args: args[2:],
	}

	return fn(cmd)
}

func (a *Application) CheckGatorCommand(args []string) (GatorCommand, error) {
	switch len(args) {
	case 0:
		return nil, NoProgramProvided
	case 1:
		if args[0] != "gator" {
			return nil, NotGatorCommandErr
		}
		return nil, NoCommandErr
	case 2:
		if args[0] != "gator" {
			return nil, NotGatorCommandErr
		}
		_, exists := a.Commands[args[1]]
		if !exists {
			return nil, NoCommandExistsErr
		}
		return nil, NoArgumentsProvidedErr
	}

	if args[0] != "gator" {
		return nil, NotGatorCommandErr
	}

	cmd, exists := a.Commands[args[1]]
	if !exists {
		return nil, NoCommandExistsErr
	}

	return cmd, nil
}

func (a *Application) Register(name string, fn func(Command) error) {
	a.Commands[name] = fn
}

func (a *Application) RegisterCommands() {
	a.Commands = make(map[string]GatorCommand)

	a.Register("login", a.loginUserHandler)
	a.Register("register", a.createUserHandler)
}

func (c *Config) ReadConfig() {
	dirPath := getGatorConfig()

	content, err := os.ReadFile(dirPath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, c)
	if err != nil {
		panic(err)
	}
}

func (c *Config) SetUser(user string) {
	c.ReadConfig()
	c.Current_user_name = user

	jsonString, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	Write(jsonString)
}

func Write(data []byte) {
	dirPath := getGatorConfig()
	err := os.WriteFile(dirPath, data, 0644)
	if err != nil {
		panic(err)
	}
}
