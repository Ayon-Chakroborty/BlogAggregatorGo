package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"ayonchakroborty.net/blogaggregator/internal/data"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	Models   data.Models
}

func main() {
	c := &Config{}
	c.ReadConfig()

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	app := &Application{
		Config: c,
		Models: data.NewModels(db),
	}

	app.RegisterCommands()

	running := true

	for running {
		err := app.Run()
		if err != nil {
			log.Fatal(err.Error())
		}
		running = false
	}
}

const dsnEnvVar string = "GATOR_DB_DSN"

func openDB() (*sql.DB, error) {
	// load DSN
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dsn := os.Getenv(dsnEnvVar)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(15 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
