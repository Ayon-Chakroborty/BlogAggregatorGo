package main

import (
	"fmt"
	"log"
	"os"
)

const (
	gatorConfig = "/.gatorconfig.json"
)

func getGatorConfig() string {
	HOMEDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return fmt.Sprint(HOMEDir + gatorConfig)
}

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

type Application struct {
	Config   *Config
	Commands map[string]GatorCommand
}

func main() {
	c := &Config{}
	c.ReadConfig()

	app := &Application{
		Config: c,
	}

	app.RegisterCommands()

	running := true

	for running {
		err := app.Run()
		if err != nil{
			log.Fatal(err.Error())
		}
		running = false
	}
}
