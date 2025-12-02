package main

import (
	"log"

	"ayonchakroborty.net/blogaggregator/internal/data"
)

func (a *Application) loginUserHandler(cmd Command) error {
	if len(cmd.Args) > 1 {
		return TooManyArgumentsErr
	}

	a.Config.SetUser(cmd.Args[0])
	log.Printf("%q has been set as the current user", cmd.Args[0])
	return nil
}

func (a *Application) createUserHandler(cmd Command) error {
	if len(cmd.Args) > 1 {
		return TooManyArgumentsErr
	}

	name := cmd.Args[0]

	user := &data.User{
		Name: name,
	}

	err := a.Models.UsersModel.Insert(user)
	if err != nil{
		log.Println(err)
	}
	log.Printf("%+v\n", *user)

	return nil
}
