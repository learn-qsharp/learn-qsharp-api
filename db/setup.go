package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"os"
)

func SetupDB(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
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
