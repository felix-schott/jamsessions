package dbutils

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxgeom "github.com/twpayne/pgx-geom"
)

func CreatePool(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("Please provide a postgres connection string using the environment variable DB_URL")
	}
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal(err)
	}

	// register geometry types
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		if err := pgxgeom.Register(ctx, conn); err != nil {
			return err
		}
		return nil
	}
	return pgxpool.NewWithConfig(context.Background(), config)
}

// func NewConn(ctx context.Context) (*sql.DB, error) {
// 	connStr := os.Getenv("DB_URL")
// 	if connStr == "" {
// 		log.Fatal("Please provide a postgres connection string using the environment variable DB_URL")
// 	}
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }
