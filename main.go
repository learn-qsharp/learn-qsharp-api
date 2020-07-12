package main

import (
	"github.com/joho/godotenv"
	"github.com/learn-qsharp/learn-qsharp-api/db"
	"github.com/learn-qsharp/learn-qsharp-api/github"
	"github.com/learn-qsharp/learn-qsharp-api/router"
	"github.com/learn-qsharp/learn-qsharp-api/tutorials"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load()

	dbc, err := db.SetupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbc.Close()

	if os.Getenv("GITHUB_IGNORE") != "true" {
		githubClient, githubCtx := github.Setup()

		err = tutorials.LoadFromGithubAndSaveToDb(dbc, githubClient, githubCtx)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = router.Run(dbc)
	if err != nil {
		log.Fatal(err)
	}
}
