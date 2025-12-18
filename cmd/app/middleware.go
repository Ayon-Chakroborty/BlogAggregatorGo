package main

import (
	"log"

	"ayonchakroborty.net/blogaggregator/internal/data"
)

func (a *Application) middlewareLoggedIn(handler func(c Command, user data.User) error) func(Command) error {
	return func(c Command) error {
		user, err := a.Models.UsersModel.Get(a.Config.Current_user_name)
		if err != nil{
			log.Fatal(err)
		}

		return handler(c, *user)
	}
}
