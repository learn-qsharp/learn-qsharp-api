package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func Setup(ctx context.Context) *github.Client {
	tc := oauth2.NewClient(ctx, nil)

	return github.NewClient(tc)
}
