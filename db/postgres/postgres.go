package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
	"github.com/jackc/pgx/v4"
)

var _ link.LinkStore = &PgDB{}

const urlExample = "postgres://postgres:postgres@localhost:5432/shortenerdb"

type PgDB struct {
	sync.Mutex
	dbConn *pgx.Conn
}

func NewPgDB() (*PgDB, error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	//conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &PgDB{
		dbConn: conn,
	}, nil
}

func (db *PgDB) CreateShortLink(ctx context.Context, longLink string) (*link.Link, error) {
	if db.dbConn == nil {
		return nil, errors.New("no connection to db")
	}
	db.Lock()
	defer db.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	sqlStatement := fmt.Sprintf("insert into links (longurl) values (%v) returning shorturl;", longLink)
	shorturl := 0

	err := db.dbConn.QueryRow(ctx, sqlStatement).Scan(&shorturl)
	if err != nil {
		return nil, fmt.Errorf("error while createShortLink, err: %w", err)
	}

	return &link.Link{
		ShortLink: shorturl,
		LongLink:  longLink,
	}, nil
}

func (db *PgDB) UpdateStat(ctx context.Context, shortLink int, ip string) error {
	if db.dbConn == nil {
		return errors.New("no connection to db")
	}
	db.Lock()
	defer db.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	sqlStatement := fmt.Sprintf("insert into stat (shorturl, userip, passtime) values (%v, %v, %v);", shortLink, ip, time.Now().String())

	//err := db.dbConn.QueryRow(ctx, sqlStatement)
	_, err := db.dbConn.Exec(ctx, sqlStatement)
	if err != nil {
		return fmt.Errorf("error while update stat table: %w", err)
	}

	return nil
}

func (db *PgDB) GetStat(ctx context.Context, shortLink int) ([]link.Stat, error) {
	if db.dbConn == nil {
		return nil, errors.New("no connection to db")
	}
	db.Lock()
	defer db.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	//select (userip, passtime) from stat where shorturl=1;
	sqlStatement := fmt.Sprintf("select (userip, passtime) from stat where shorturl=%v;", shortLink)

	rows, err := db.dbConn.Query(ctx, sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error while getstat err: %w", err)
	}
	defer rows.Close()

	var stats []link.Stat

	for rows.Next() {
		var stat link.Stat
		if err := rows.Scan(&stat.UserIP, &stat.PassTime); err != nil {
			return stats, err
		}
		stats = append(stats, stat)
	}

	if err = rows.Err(); err != nil {
		return stats, err
	}

	return stats, nil
}
