package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/learn-qsharp/learn-qsharp-api/db"
	"github.com/learn-qsharp/learn-qsharp-api/github"
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
	defer dbc.Close(context.Background())

	if os.Getenv("GITHUB_IGNORE") != "true" {
		ctx := context.Background()

		githubClient := github.Setup(ctx)

		err = tutorials.LoadFromGithubAndSaveToDb(ctx, dbc, githubClient)
		if err != nil {
			log.Fatal(err)
		}
	}

	/*err = router.Run(dbc)
	if err != nil {
		log.Fatal(err)
	}*/
}
