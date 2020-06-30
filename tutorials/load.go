package tutorials

import (
	"context"
	"github.com/google/go-github/v32/github"
	"log"
)

func Load(client *github.Client, ctx context.Context) error {
	ids, err := getTutorialIDs(client, ctx)
	if err != nil {
		return err
	}

	log.Println(ids)

	return nil
}

func getTutorialIDs(client *github.Client, ctx context.Context) ([]string, error) {
	opts := github.RepositoryContentGetOptions{Ref: "master"}
	_, directories, _, err := client.Repositories.GetContents(ctx, "learn-qsharp", "learn-qsharp-tutorials", "tutorials", &opts)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0)
	for _, directory := range directories {
		ids = append(ids, directory.GetName())
	}

	return ids, nil
}
