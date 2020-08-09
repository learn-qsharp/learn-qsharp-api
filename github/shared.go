package github

import (
	"context"
	"github.com/google/go-github/v32/github"
)

func getLatestBranchSHA(ctx context.Context, client *github.Client, owner, repo, branchName string) (string, error) {
	branch, _, err := client.Repositories.GetBranch(ctx, owner, repo, branchName)
	if err != nil {
		return "", err
	}
	return branch.Commit.GetSHA(), nil
}
