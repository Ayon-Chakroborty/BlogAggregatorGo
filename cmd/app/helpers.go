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
	if len(args) == 0 {
		return nil, ErrNoProgramProvided
	}

	if args[0] != "gator" {
		return nil, ErrNotGatorCommand
	}

	if len(args) < 2 {
		return nil, ErrNoCommand
	}

	cmd, exists := a.Commands[args[1]]
	if !exists {
		return nil, ErrNoCommandExists
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
	a.Register("users", a.listUsersHandler)
	a.Register("agg", a.FetchFeedHandler)
	a.Register("addFeed", a.middlewareLoggedIn(a.createFeedHandler))
	a.Register("feeds", a.getAllFeedsHandler)
	a.Register("follow", a.middlewareLoggedIn(a.FollowHanlder))
	a.Register("following", a.middlewareLoggedIn(a.FollowingHandler))
	a.Register("unfollow", a.middlewareLoggedIn(a.UnfollowHandler))
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
