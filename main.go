package main

import (
	"context"
	"github.com/learn-qsharp/learn-qsharp-api/db"
	"github.com/learn-qsharp/learn-qsharp-api/env"
	"github.com/learn-qsharp/learn-qsharp-api/github"
	"github.com/learn-qsharp/learn-qsharp-api/router"
	"log"
)

func main() {
	envVars, err := env.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pgxConn, err := db.SetupPgxConn(ctx, envVars.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pgxConn.Close(ctx)

	pgxPool, err := db.SetupPgxPool(ctx, envVars.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pgxPool.Close()

	err = db.Migrate(ctx, pgxConn)
	if err != nil {
		log.Fatal(err)
	}

	githubClient := github.Setup(ctx, envVars)

	err = github.UpdateTutorials(ctx, envVars, pgxConn, githubClient)
	if err != nil {
		log.Fatal(err)
	}

	err = github.UpdateProblems(ctx, envVars, pgxConn, githubClient)
	if err != nil {
		log.Fatal(err)
	}

	err = router.Run(pgxPool)
	if err != nil {
		log.Fatal(err)
	}
}
