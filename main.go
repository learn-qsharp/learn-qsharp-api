package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/learn-qsharp/learn-qsharp-api/db"
	"log"
)

func main() {
	_ = godotenv.Load()

	dbc, err := db.SetupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbc.Close(context.Background())

	/*if os.Getenv("GITHUB_IGNORE") != "true" {
		githubClient, githubCtx := github.Setup()

		err = tutorials.LoadFromGithubAndSaveToDb(dbc, githubClient, githubCtx)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = router.Run(dbc)
	if err != nil {
		log.Fatal(err)
	}*/
}
