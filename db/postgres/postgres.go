// Postgres package implements methods for interactions with the Postgres database.
package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

var _ link.LinkStore = &PgDB{}

const UrlExample = "postgres://postgres:postgres@localhost:5432/shortenerdb"

type PgDB struct {
	sync.Mutex
	dbConn *pgx.Conn
}

func WaitDB() {
	UrlExample := os.Getenv("PG_DSN")
	fmt.Println("--- waiting for db ---")

	for {
		conn, err := pgx.Connect(context.Background(), UrlExample)
		if err != nil {
			fmt.Println("waiting for db...")
			time.Sleep(1 * time.Second)
			continue
		}
		conn.Close(context.Background())
		break
	}
	fmt.Println("db ok")
}

// NewPgDB returns new database object (PgDB) with the installed connection to it.
func NewPgDB() (*PgDB, error) {
	WaitDB()

	UrlExample := os.Getenv("PG_DSN")
	conn, err := pgx.Connect(context.Background(), UrlExample)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	fmt.Println("Db connected.")

	return &PgDB{
		dbConn: conn,
	}, nil
}

func (db *PgDB) Close() {
	db.dbConn.Close(context.Background())
}

// CreateShortLink returns a Short URL for the provided Lonk URL.
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

	sqlStatement := fmt.Sprintf("insert into links (longurl) values ('%v') returning shorturl;", longLink)
	shorturl := 0

	err := db.dbConn.QueryRow(ctx, sqlStatement).Scan(&shorturl)
	if err != nil {
		return nil, fmt.Errorf("error while createShortLink, err: %w", err)
	}

	return &link.Link{
		ShortLink: strconv.Itoa(shorturl),
		LongLink:  longLink,
	}, nil
}

// UpdateStat updates statistics for the provided Short URL and IP address.
func (db *PgDB) UpdateStat(ctx context.Context, shortLink int, ip string) error {
	if db.dbConn == nil {
		return errors.New("no connection to db")
	}
	db.Lock()
	defer db.Unlock()

	sqlStatement := fmt.Sprintf("insert into stat (shorturl, userip, passtime) values (%v, '%v', '%v');", shortLink, ip, time.Now().String())

	_, err := db.dbConn.Exec(ctx, sqlStatement)
	if err != nil {
		return fmt.Errorf("error while update stat table: %w", err)
	}

	return nil
}

// GetStat returns statistics for the provided Short URL.
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

	sqlStatement := fmt.Sprintf("select (userip, passtime) from stat where shorturl=%v;", shortLink)

	rows, err := db.dbConn.Query(ctx, sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error while getstat err: %w", err)
	}
	defer rows.Close()

	var stats []link.Stat

	for rows.Next() {

		select {
		case <-ctx.Done():
			return stats, ctx.Err()
		default:
		}

		var stat link.Stat
		var i interface{}
		if err := rows.Scan(&i); err != nil {
			return stats, err
		}

		values, ok := i.([]pgtype.Value)
		if !ok {
			return stats, errors.New("wrong type returned from db")
		}

		if len(values) != 2 {
			return stats, errors.New("wrong columns number returned from db")
		}

		stat.UserIP, ok = values[0].Get().(string)
		if !ok {
			return stats, errors.New("wrong type of userIP returned from db")
		}

		stat.PassTime, ok = values[1].Get().(string)
		if !ok {
			return stats, errors.New("wrong type of PassTime returned from db")
		}

		stats = append(stats, stat)
	}

	if err = rows.Err(); err != nil {
		return stats, err
	}

	return stats, nil
}

// GetLongLink returns Long URL for the provided Short URL.
func (db *PgDB) GetLongLink(ctx context.Context, shortLink int) (link.Link, error) {
	if db.dbConn == nil {
		return link.Link{}, errors.New("no connection to db")
	}
	db.Lock()
	defer db.Unlock()

	sqlStatement := fmt.Sprintf("select (longurl) from links where shorturl='%v';", shortLink)

	lurl := ""

	err := db.dbConn.QueryRow(ctx, sqlStatement).Scan(&lurl)
	if err != nil {
		return link.Link{}, fmt.Errorf("error while GetLongLink, err: %w", err)
	}

	shortLinkStr := strconv.Itoa(shortLink)

	return link.Link{LongLink: lurl, ShortLink: shortLinkStr}, nil
}
