package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func Setup() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "... your access token ..."},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
