package main

import (
	"github.com/joho/godotenv"
	"github.com/learn-qsharp/learn-qsharp-api/db"
	"github.com/learn-qsharp/learn-qsharp-api/github"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	dbc, err := db.SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	_ = github.Setup()

	defer dbc.Close()
}
