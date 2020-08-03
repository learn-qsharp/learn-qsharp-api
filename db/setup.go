package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

func SetupPgxConn(ctx context.Context, databaseURL string) (*pgx.Conn, error) {
	return pgx.Connect(ctx, databaseURL)
}

func SetupPgxPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, databaseURL)
}

func Migrate(ctx context.Context, pgxConn *pgx.Conn) error {
	m, err := migrate.NewMigrator(ctx, pgxConn, "schema_version")
	if err != nil {
		return err
	}

	err = m.LoadMigrations("migrations")
	if err != nil {
		return err
	}

	return m.Migrate(ctx)
}
