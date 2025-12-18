package main

import (
	"fmt"
	"log"

	"ayonchakroborty.net/blogaggregator/internal/data"
)

func (a *Application) FollowHanlder(cmd Command, user data.User) error {
	if len(cmd.Args) < 1 {
		return ErrNoArgumentsProvided
	}

	if len(cmd.Args) > 1 {
		return ErrTooManyArguments
	}

	url := cmd.Args[0]

	feed, err := a.Models.FeedsModel.Get(url)
	if err != nil {
		log.Fatal(err.Error())
	}

	record := &data.FeedFollow{UserId: user.Id, FeedId: feed.Id}

	err = a.Models.FeedFollowsModel.Insert(record)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%q is now follwoing feed %q\n", record.UserName, record.FeedName)

	return nil
}

func (a *Application) FollowingHandler(cmd Command, user data.User) error {
	if len(cmd.Args) > 0 {
		return ErrTooManyArguments
	}

	feeds, err := a.Models.FeedFollowsModel.GetAll(user.Id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%q is following these feeds:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("feed name: %q, url: %q\n", feed.FeedName, feed.Url)
	}

	return nil
}

func (a *Application) UnfollowHandler(cmd Command, user data.User) error {
	if len(cmd.Args) < 1{
		return ErrNoArgumentsProvided
	}

	if len(cmd.Args) > 1 {
		return ErrTooManyArguments
	}

	url := cmd.Args[0]

	feed, err := a.Models.FeedsModel.Get(url)
	if err != nil{
		log.Fatal(err)
	}

	err = a.Models.FeedFollowsModel.Delete(feed.Id, user.Id)
	if err != nil{
		log.Fatal(err)
	}

	fmt.Printf("%q has unfollowed %q\n", user.Name, feed.Url)
	
	return nil
}
