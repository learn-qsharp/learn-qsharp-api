package github

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
	"github.com/learn-qsharp/learn-qsharp-api/env"
	"github.com/learn-qsharp/learn-qsharp-api/models"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"strconv"
	"strings"
)

type tutorialMetadata struct {
	Title       string
	Credits     string
	Description string
	Difficulty  string
	Tags        []string
}

func UpdateTutorials(ctx context.Context, envVars env.Env, db *pgx.Conn, client *github.Client) error {
	hash, err := getLatestBranchSHA(ctx, envVars, client)
	if err != nil {
		return err
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	mustBeUpdated, err := mustBeUpdated(ctx, tx, hash)
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

func getLatestBranchSHA(ctx context.Context, envVars env.Env, client *github.Client) (string, error) {
	branch, _, err := client.Repositories.GetBranch(ctx, envVars.GithubTutorialsOwner, envVars.GithubTutorialsRepo,
		envVars.GithubTutorialsBranch)
	if err != nil {
		return "", err
	}
	return branch.Commit.GetSHA(), nil
}

func mustBeUpdated(ctx context.Context, tx pgx.Tx, githubHash string) (bool, error) {
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

func upsertLatestHash(ctx context.Context, tx pgx.Tx, githubHash string) error {
	sql := `
		INSERT INTO tutorials_hash
		VALUES(1, $1)
		ON CONFLICT (id)
		DO
			UPDATE SET hash = $1
	`

	_, err := tx.Exec(ctx, sql, githubHash)

	return err
}

func loadTutorials(ctx context.Context, envVars env.Env, client *github.Client) ([]models.Tutorial, error) {
	ids, err := getTutorialIDs(ctx, envVars, client)
	if err != nil {
		return nil, err
	}

	var tutorials []models.Tutorial
	for _, id := range ids {
		body, err := getTutorialBody(ctx, envVars, client, id)
		if err != nil {
			return nil, err
		}

		metadata, err := getTutorialMetadata(ctx, envVars, client, id)
		if err != nil {
			return nil, err
		}

		tutorial := models.Tutorial{
			ID:          id,
			Title:       metadata.Title,
			Credits:     metadata.Credits,
			Description: metadata.Description,
			Body:        body,
			Difficulty:  metadata.Difficulty,
			Tags:        metadata.Tags,
		}

		tutorials = append(tutorials, tutorial)
	}

	return tutorials, nil
}

func getTutorialIDs(ctx context.Context, envVars env.Env, client *github.Client) ([]uint, error) {
	opts := github.RepositoryContentGetOptions{Ref: envVars.GithubTutorialsBranch}
	_, directories, _, err := client.Repositories.GetContents(
		ctx,
		envVars.GithubTutorialsOwner,
		envVars.GithubTutorialsRepo,
		"tutorials", &opts,
	)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, 0)
	for _, directory := range directories {
		id, err := strconv.Atoi(directory.GetName())
		if err != nil {
			return nil, err
		}

		if id <= 0 {
			return nil, errors.New("id must be positive")
		}

		ids = append(ids, uint(id))
	}

	return ids, nil
}

func getTutorialBody(ctx context.Context, envVars env.Env, client *github.Client, id uint) (string, error) {
	opts := github.RepositoryContentGetOptions{Ref: envVars.GithubTutorialsBranch}
	r, err := client.Repositories.DownloadContents(
		ctx,
		envVars.GithubTutorialsOwner,
		envVars.GithubTutorialsRepo,
		fmt.Sprintf("tutorials/%d/body.md", id),
		&opts,
	)
	if err != nil {
		return "", err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getTutorialMetadata(ctx context.Context, envVars env.Env, client *github.Client, id uint) (*tutorialMetadata, error) {
	opts := github.RepositoryContentGetOptions{Ref: envVars.GithubTutorialsBranch}
	r, err := client.Repositories.DownloadContents(
		ctx,
		envVars.GithubTutorialsOwner,
		envVars.GithubTutorialsRepo,
		fmt.Sprintf("tutorials/%d/metadata.yaml", id),
		&opts,
	)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	metadata := tutorialMetadata{}
	err = yaml.Unmarshal(buf.Bytes(), &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func createOrUpdateTutorialsOnDatabase(ctx context.Context, tx pgx.Tx, tutorials []models.Tutorial) error {
	for _, tutorial := range tutorials {
		if err := createOrUpdateTutorialOnDatabase(ctx, tx, &tutorial); err != nil {
			return err
		}
	}

	return nil
}

func createOrUpdateTutorialOnDatabase(ctx context.Context, tx pgx.Tx, tutorial *models.Tutorial) error {
	if tutorial == nil {
		return errors.New("tutorial can't be nil")
	}

	sql := `
		INSERT INTO tutorials (id, title, credits, description, body, difficulty, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			title = $2,
			credits = $3,
			description = $4,
			body = $5,
			difficulty = $6,
			tags = $7
`

	_, err := tx.Exec(ctx, sql, tutorial.ID, tutorial.Title, tutorial.Credits, tutorial.Description, tutorial.Body,
		tutorial.Difficulty, tutorial.Tags)

	return err
}
