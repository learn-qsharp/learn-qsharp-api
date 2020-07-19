package tutorials

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
	"github.com/learn-qsharp/learn-qsharp-api/models"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strconv"
	"strings"
)

type metadata struct {
	Title       string
	Credits     string
	Description string
	Difficulty  string
	Tags        []string
}

func LoadFromGithubAndSaveToDb(ctx context.Context, db *pgx.Conn, client *github.Client) error {
	tutorials, err := loadTutorials(client, ctx)
	if err != nil {
		return err
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err = createOrUpdateTutorialsOnDatabase(ctx, tx, tutorials); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func loadTutorials(client *github.Client, ctx context.Context) ([]models.Tutorial, error) {
	ids, err := getTutorialIDs(ctx, client)
	if err != nil {
		return nil, err
	}

	var tutorials []models.Tutorial
	for _, id := range ids {
		body, err := getTutorialBody(ctx, client, id)
		if err != nil {
			return nil, err
		}

		metadata, err := getTutorialMetadata(ctx, client, id)
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

func getTutorialIDs(ctx context.Context, client *github.Client) ([]uint, error) {
	opts := github.RepositoryContentGetOptions{Ref: os.Getenv("GITHUB_TUTORIALS_REF")}
	_, directories, _, err := client.Repositories.GetContents(
		ctx,
		os.Getenv("GITHUB_TUTORIALS_OWNER"),
		os.Getenv("GITHUB_TUTORIALS_REPO"),
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

func getTutorialBody(ctx context.Context, client *github.Client, id uint) (string, error) {
	opts := github.RepositoryContentGetOptions{Ref: os.Getenv("GITHUB_TUTORIALS_REF")}
	r, err := client.Repositories.DownloadContents(
		ctx,
		os.Getenv("GITHUB_TUTORIALS_OWNER"),
		os.Getenv("GITHUB_TUTORIALS_REPO"),
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

func getTutorialMetadata(ctx context.Context, client *github.Client, id uint) (*metadata, error) {
	opts := github.RepositoryContentGetOptions{Ref: os.Getenv("GITHUB_TUTORIALS_REF")}
	r, err := client.Repositories.DownloadContents(
		ctx,
		os.Getenv("GITHUB_TUTORIALS_OWNER"),
		os.Getenv("GITHUB_TUTORIALS_REPO"),
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

	metadata := metadata{}
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

	/*searchTutorial := &models.Tutorial{
		ID: tutorial.ID,
	}*/

	/*if err := tx.First(searchTutorial).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err = tx.Create(tutorial).Error; err != nil {
				return err
			}

			return nil
		} else {
			return err
		}
	}

	// It will update only changed fields.
	if err := tx.Model(searchTutorial).Updates(tutorial).Error; err != nil {
		return err
	}*/

	return nil
}
