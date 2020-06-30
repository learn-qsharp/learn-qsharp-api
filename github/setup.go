package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"os"
)

func Setup() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_AUTH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), ctx
}
