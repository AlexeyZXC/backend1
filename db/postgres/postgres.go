package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
	"github.com/jackc/pgx/v4"
)

var _ link.LinkStore = &PgDB{}

type PgDB struct {
	dbConn *pgx.Conn
}

func NewPgDB() (*PgDB, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	return &PgDB{
		dbConn: conn,
	}, nil
}

func (db *PgDB) CreateShortLink(ctx context.Context, longLink string) (*link.Link, error) {
	return nil, nil
}

func (db *PgDB) UpdateStat(ctx context.Context, shortLink int, ip string) error {
	return nil
}

func (db *PgDB) GetStat(ctx context.Context, shortLink int) ([]link.Stat, error) {
	return nil, nil
}
