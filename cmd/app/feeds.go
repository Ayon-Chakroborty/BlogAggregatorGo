package main

import (
	"log"

	"ayonchakroborty.net/blogaggregator/internal/data"
)

func (a *Application) createFeedHandler(cmd Command, user data.User) error {
	if len(cmd.Args) < 2 {
		return ErrNotEnoughArguments
	}

	if len(cmd.Args) > 2 {
		return ErrTooManyArguments
	}

	feedName := cmd.Args[0]
	url := cmd.Args[1]

	feed := data.Feed{
		Name:   feedName,
		Url:    url,
		UserId: user.Id,
	}

	err := a.Models.FeedsModel.Insert(&feed)
	if err != nil {
		log.Fatal(err.Error())
	}

	feedFollow := data.FeedFollow{UserId: feed.UserId, FeedId: feed.Id}
	err = a.Models.FeedFollowsModel.Insert(&feedFollow)
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
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, feed := range feeds {
		log.Printf("%+v\n", feed)
	}

	return nil
}
