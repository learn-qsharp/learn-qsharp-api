package env

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseURL string `env:"DATABASE_URL"`

	GithubIgnore         bool   `env:"GITHUB_IGNORE"`
	GithubTutorialsOwner string `env:"GITHUB_TUTORIALS_OWNER"`
	GithubTutorialsRepo  string `env:"GITHUB_TUTORIALS_REPO"`
	GithubTutorialsRef   string `env:"GITHUB_TUTORIALS_REF"`
}

func Load() (Env, error) {
	_ = godotenv.Load()

	envVars := Env{}
	if err := env.Parse(&envVars); err != nil {
		return Env{}, err
	}

	return envVars, nil
}
