package main

import (
	"fmt"
	"log"

	"ayonchakroborty.net/blogaggregator/internal/data"
)

func (a *Application) loginUserHandler(cmd Command) error {
	if len(cmd.Args) < 1 {
		return ErrNoArgumentsProvided
	}

	if len(cmd.Args) > 1 {
		return ErrTooManyArguments
	}

	// Check if user is registered
	name := cmd.Args[0]

	user, err := a.Models.UsersModel.Get(name)
	if err != nil {
		return err
	}

	a.Config.SetUser(user.Name)
	log.Printf("%q has been set as the current user", cmd.Args[0])
	return nil
}

func (a *Application) createUserHandler(cmd Command) error {
	if len(cmd.Args) < 1 {
		return ErrNoArgumentsProvided
	}
	
	if len(cmd.Args) > 1 {
		return ErrTooManyArguments
	}

	name := cmd.Args[0]

	user := &data.User{
		Name: name,
	}

	err := a.Models.UsersModel.Insert(user)
	if err != nil {
		return err
	}

	a.Config.SetUser(user.Name)
	return nil
}

func (a *Application) listUsersHandler(cmd Command) error {
	if len(cmd.Args) > 0 {
		return ErrTooManyArguments
	}

	users, err := a.Models.UsersModel.ListUsers()
	if err != nil {
		return err
	}

	currentUsr := a.Config.Current_user_name

	for _, user := range users {
		fmt.Print("* ")
		fmt.Printf("%s", user.Name)
		if currentUsr == user.Name {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}
