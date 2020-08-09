package github

import (
	"context"
	"errors"
	"github.com/google/go-github/v32/github"
	"strconv"
)

func getLatestBranchSHA(ctx context.Context, client *github.Client, owner, repo, branchName string) (string, error) {
	branch, _, err := client.Repositories.GetBranch(ctx, owner, repo, branchName)
	if err != nil {
		return "", err
	}
	return branch.Commit.GetSHA(), nil
}

func getIDsFromDirectories(directories []*github.RepositoryContent) ([]uint, error) {
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
