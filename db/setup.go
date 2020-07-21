package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/learn-qsharp/learn-qsharp-api/env"
)

func SetupDB(ctx context.Context, envVars env.Env) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, envVars.DatabaseURL)
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		return nil, err
	}

	err = m.LoadMigrations("migrations")
	if err != nil {
		return nil, err
	}

	err = m.Migrate(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
