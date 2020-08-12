package env

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseURL string `env:"DATABASE_URL"`

	GithubBotToken string `env:"GITHUB_BOT_TOKEN"`

	GithubTutorialsOwner  string `env:"GITHUB_TUTORIALS_OWNER"`
	GithubTutorialsRepo   string `env:"GITHUB_TUTORIALS_REPO"`
	GithubTutorialsBranch string `env:"GITHUB_TUTORIALS_BRANCH"`

	GithubProblemsOwner  string `env:"GITHUB_PROBLEMS_OWNER"`
	GithubProblemsRepo   string `env:"GITHUB_PROBLEMS_REPO"`
	GithubProblemsBranch string `env:"GITHUB_PROBLEMS_BRANCH"`
}

func Load() (Env, error) {
	_ = godotenv.Load()

	envVars := Env{}
	if err := env.Parse(&envVars); err != nil {
		return Env{}, err
	}

	return envVars, nil
}
