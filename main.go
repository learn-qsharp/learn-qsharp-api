package main

import (
	"context"
	"github.com/learn-qsharp/learn-qsharp-api/db"
	"github.com/learn-qsharp/learn-qsharp-api/env"
	"github.com/learn-qsharp/learn-qsharp-api/github"
	"github.com/learn-qsharp/learn-qsharp-api/router"
	"github.com/learn-qsharp/learn-qsharp-api/tutorials"
	"log"
)

func main() {
	envVars, err := env.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	dbc, err := db.SetupDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer dbc.Close(ctx)

	if !envVars.GithubIgnore {
		githubClient := github.Setup(ctx)

		err = tutorials.LoadFromGithubAndSaveToDb(ctx, dbc, githubClient)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = router.Run(dbc)
	if err != nil {
		log.Fatal(err)
	}
}
