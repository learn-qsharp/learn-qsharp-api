package tutorials

import (
	"bytes"
	"context"
	"github.com/google/go-github/v32/github"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"strings"
)

type metadata struct {
	Title      string
	Author     string
	Difficulty string
	Tags       []string
}

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

		metadata, err := getTutorialMetadata(id, client, ctx)
		if err != nil {
			return err
		}

		log.Println(description)
		log.Println(metadata)
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

func getTutorialMetadata(id string, client *github.Client, ctx context.Context) (*metadata, error) {
	opts := github.RepositoryContentGetOptions{Ref: "master"}
	r, err := client.Repositories.DownloadContents(
		ctx,
		os.Getenv("GITHUB_TUTORIALS_OWNER"),
		os.Getenv("GITHUB_TUTORIALS_REPO"),
		"tutorials/"+id+"/metadata.yaml",
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
