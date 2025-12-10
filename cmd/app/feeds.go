package main

import (
	"log"

	"ayonchakroborty.net/blogaggregator/internal/data"
)

func (a *Application) createFeedHandler(cmd Command) error {
	if len(cmd.Args) < 2 {
		return ErrNotEnoughArguments
	}

	if len(cmd.Args) > 2 {
		return ErrTooManyArguments
	}

	name := a.Config.Current_user_name

	user, err := a.Models.UsersModel.Get(name)
	if err != nil{
		return err
	}

	feedName := cmd.Args[0]
	url := cmd.Args[1]

	feed := data.Feed{
		Name:    feedName,
		Url:     url,
		User_id: user.Id,
	}

	err = a.Models.FeedsModel.Insert(&feed)
	if err != nil {
		log.Fatal(err.Error())
	}
	
	log.Printf("Successfully added feed: %+v\n", feed)
	
	return nil
}

func (a *Application) getAllFeedsHandler(cmd Command) error {
	if len(cmd.Args) > 0 {
		return ErrTooManyArguments
	}

	feeds, err := a.Models.FeedsModel.GetAll()
	if err != nil{
		log.Fatal(err.Error())
	}

	for _, feed  := range feeds{
		log.Printf("%+v\n", feed)
	}

	return nil
}
