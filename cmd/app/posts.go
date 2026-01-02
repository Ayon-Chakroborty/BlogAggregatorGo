package main

import (
	"fmt"
	"log"
	"strconv"

	"ayonchakroborty.net/blogaggregator/internal/data"
)

func (a *Application) BrowseHandler(cmd Command, user data.User) error {
	if len(cmd.Args) > 1 {
		return ErrTooManyArguments
	}

	if len(cmd.Args) == 0 {
		cmd.Args = append(cmd.Args, "0")
	}

	limit, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return err
	}

	posts, err := a.Models.PostsModel.GetAll(&user, limit)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		fmt.Printf("%+v\n", post)
	}

	return nil
}
