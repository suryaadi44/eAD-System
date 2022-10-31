package main

import (
	"github.com/joho/godotenv"
	"github.com/suryaadi44/eAD-System/pkg/config"
	"log"

	"github.com/suryaadi44/eAD-System/pkg/database"
)

func init() {
	//	load env variables from .env file for local development
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	env := config.LoadConfig()

	db, err := database.Connect(
		env["DB_HOST"],
		env["DB_PORT"],
		env["DB_USER"],
		env["DB_PASS"],
		env["DB_NAME"],
		5,
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
