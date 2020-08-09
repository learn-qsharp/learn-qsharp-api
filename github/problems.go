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
	"strings"
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
		log.Println("Table problems is up-to-date.")
		return nil
	}

	problems, err := loadProblems(ctx, envVars, client)
	if err != nil {
		return err
	}

	if err = createOrUpdateProblemsOnDatabase(ctx, tx, problems); err != nil {
		return err
	}

	err = upsertLatestProblemsHash(ctx, tx, hash)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func mustProblemsBeUpdated(ctx context.Context, tx pgx.Tx, githubHash string) (bool, error) {
	var dbHash string
	err := tx.QueryRow(ctx, "SELECT hash FROM problems_hash WHERE id = 1").Scan(&dbHash)
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

func upsertLatestProblemsHash(ctx context.Context, tx pgx.Tx, githubHash string) error {
	sql := `
		INSERT INTO problems_hash
		VALUES(1, $1)
		ON CONFLICT (id)
		DO
			UPDATE SET hash = $1
	`

	_, err := tx.Exec(ctx, sql, githubHash)

	return err
}

func loadProblems(ctx context.Context, envVars env.Env, client *github.Client) ([]models.Problem, error) {
	ids, err := getProblemIDs(ctx, envVars, client)
	if err != nil {
		return nil, err
	}

	var problems []models.Problem
	for _, id := range ids {
		body, err := getProblemBody(ctx, envVars, client, id)
		if err != nil {
			return nil, err
		}

		template, err := getProblemTemplate(ctx, envVars, client, id)
		if err != nil {
			return nil, err
		}

		metadata, err := getProblemMetadata(ctx, envVars, client, id)
		if err != nil {
			return nil, err
		}

		problem := models.Problem{
			ID:         id,
			Name:       metadata.Name,
			Credits:    metadata.Credits,
			Body:       body,
			Template:   template,
			Difficulty: metadata.Difficulty,
			Tags:       metadata.Tags,
		}

		problems = append(problems, problem)
	}

	return problems, nil
}

func getProblemIDs(ctx context.Context, envVars env.Env, client *github.Client) ([]uint, error) {
	opts := github.RepositoryContentGetOptions{Ref: envVars.GithubProblemsBranch}
	_, directories, _, err := client.Repositories.GetContents(
		ctx,
		envVars.GithubProblemsOwner,
		envVars.GithubProblemsRepo,
		"problems", &opts,
	)
	if err != nil {
		return nil, err
	}

	return getIDsFromDirectories(directories)
}

func getProblemBody(ctx context.Context, envVars env.Env, client *github.Client, id uint) (string, error) {
	opts := github.RepositoryContentGetOptions{Ref: envVars.GithubProblemsBranch}
	r, err := client.Repositories.DownloadContents(
		ctx,
		envVars.GithubProblemsOwner,
		envVars.GithubProblemsRepo,
		fmt.Sprintf("problems/%d/body.md", id),
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

func getProblemTemplate(ctx context.Context, envVars env.Env, client *github.Client, id uint) (string, error) {
	opts := github.RepositoryContentGetOptions{Ref: envVars.GithubProblemsBranch}
	r, err := client.Repositories.DownloadContents(
		ctx,
		envVars.GithubProblemsOwner,
		envVars.GithubProblemsRepo,
		fmt.Sprintf("problems/%d/template.qs", id),
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

func getProblemMetadata(ctx context.Context, envVars env.Env, client *github.Client, id uint) (*problemMetadata, error) {
	opts := github.RepositoryContentGetOptions{Ref: envVars.GithubProblemsBranch}
	r, err := client.Repositories.DownloadContents(
		ctx,
		envVars.GithubProblemsOwner,
		envVars.GithubProblemsRepo,
		fmt.Sprintf("problems/%d/metadata.yaml", id),
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

	metadata := problemMetadata{}
	err = yaml.Unmarshal(buf.Bytes(), &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func createOrUpdateProblemsOnDatabase(ctx context.Context, tx pgx.Tx, problems []models.Problem) error {
	for _, problem := range problems {
		if err := createOrUpdateProblemOnDatabase(ctx, tx, &problem); err != nil {
			return err
		}
	}

	return nil
}

func createOrUpdateProblemOnDatabase(ctx context.Context, tx pgx.Tx, problem *models.Problem) error {
	if problem == nil {
		return errors.New("problem can't be nil")
	}

	sql := `
		INSERT INTO problems (id, name, credits, body, template, difficulty, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = $2,
			credits = $3,
			body = $4,
			template = $5,
			difficulty = $6,
			tags = $7
`

	_, err := tx.Exec(ctx, sql, problem.ID, problem.Name, problem.Credits, problem.Body, problem.Template,
		problem.Difficulty, problem.Tags)

	return err
}
