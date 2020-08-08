package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/learn-qsharp/learn-qsharp-api/env"
	"golang.org/x/oauth2"
)

func Setup(ctx context.Context, envVars env.Env) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: envVars.GithubBotToken},
	)

	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
