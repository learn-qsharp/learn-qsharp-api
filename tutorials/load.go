package tutorials

import (
	"context"
	"github.com/google/go-github/v32/github"
	"io"
	"log"
	"os"
	"strings"
)

func Load(client *github.Client, ctx context.Context) error {
	ids, err := getTutorialIDs(client, ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {
		description, err := getTutorialDescription(id, client, ctx)
		if err != nil {
			return err
		}

		log.Println(description)
	}

	return nil
}

func getTutorialIDs(client *github.Client, ctx context.Context) ([]string, error) {
	opts := github.RepositoryContentGetOptions{Ref: "master"}
	_, directories, _, err := client.Repositories.GetContents(
		ctx,
		os.Getenv("GITHUB_TUTORIALS_OWNER"),
		os.Getenv("GITHUB_TUTORIALS_REPO"),
		"tutorials", &opts,
	)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0)
	for _, directory := range directories {
		ids = append(ids, directory.GetName())
	}

	return ids, nil
}

func getTutorialDescription(id string, client *github.Client, ctx context.Context) (string, error) {
	opts := github.RepositoryContentGetOptions{Ref: "master"}
	r, err := client.Repositories.DownloadContents(
		ctx,
		os.Getenv("GITHUB_TUTORIALS_OWNER"),
		os.Getenv("GITHUB_TUTORIALS_REPO"),
		"tutorials/"+id+"/description.md",
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
