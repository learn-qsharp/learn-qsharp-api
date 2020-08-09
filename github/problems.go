package github

import (
	"context"
	"errors"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
	"github.com/learn-qsharp/learn-qsharp-api/env"
	"log"
)

type problemMetadata struct {
	Name       string
	Credits    string
	Difficulty string
	Tags       []string
}

func UpdateProblems(ctx context.Context, envVars env.Env, db *pgx.Conn, client *github.Client) error {
	hash, err := getLatestBranchSHA(ctx, client, envVars.GithubProblemsOwner, envVars.GithubProblemsRepo,
		envVars.GithubProblemsBranch)
	if err != nil {
		return err
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	mustBeUpdated, err := mustProblemsBeUpdated(ctx, tx, hash)
	if err != nil {
		return err
	}

	if !mustBeUpdated {
		log.Println("Table tutorials is up-to-date.")
		return nil
	}

	tutorials, err := loadTutorials(ctx, envVars, client)
	if err != nil {
		return err
	}

	if err = createOrUpdateTutorialsOnDatabase(ctx, tx, tutorials); err != nil {
		return err
	}

	err = upsertLatestHash(ctx, tx, hash)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func mustProblemsBeUpdated(ctx context.Context, tx pgx.Tx, githubHash string) (bool, error) {
	var dbHash string
	err := tx.QueryRow(ctx, "SELECT hash FROM tutorials_hash WHERE id = 1").Scan(&dbHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		} else {
			return true, err
		}
	}

	if githubHash != dbHash {
		return true, nil
	}

	return false, nil
}
