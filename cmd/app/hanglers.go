package main

import (
	"log"
)

func (a *Application) handlerLogin(cmd Command) error {
	if len(cmd.Args) == 0 {
		return NoArgumentsProvidedErr
	}

	if len(cmd.Args) > 1 {
		return TooManyArgumentsErr
	}

	a.Config.SetUser(cmd.Args[0])
	log.Printf("%q has been set as the current user", cmd.Args[0])
	return nil
}
