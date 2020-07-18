package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"os"
)

func SetupDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
