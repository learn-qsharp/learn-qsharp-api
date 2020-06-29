package main

import (
	"github.com/joho/godotenv"
	"github.com/learn-qsharp/learn-qsharp-api/db"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbc, err := db.SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	defer dbc.Close()
}
