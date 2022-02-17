package db

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"

	"log"
	"os"

	"github.com/juanjcsr/twittlks/lks"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type DBClient struct {
	db    *sql.DB
	BunDB bun.DB
}

func OpenSQLConn(url string) *DBClient {
	dsn := url

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	return &DBClient{
		db: sqldb,
	}
}

func (db *DBClient) OpenBUN() {
	d := pgdialect.New()
	database := bun.NewDB(db.db, d)
	database.AddQueryHook(bundebug.NewQueryHook())
	db.BunDB = *database
}

func (db *DBClient) CreateTables(ctx context.Context, refresh bool) error {
	if refresh {
		err := db.BunDB.ResetModel(ctx, (*lks.Users)(nil),
			(*lks.TuitLike)(nil), (*lks.Media)(nil))
		if err != nil {
			return err
		}
	}
	_, err := db.BunDB.NewCreateTable().
		Model((*lks.Users)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.BunDB.NewCreateTable().
		Model((*lks.TuitLike)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.BunDB.NewCreateTable().
		Model((*lks.Media)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBClient) GetLastInsertedTuit(ctx context.Context) (string, error) {
	r, err := d.BunDB.QueryContext(ctx, "select id from tuit_likes order by inserted_at desc limit 1")
	if err != nil {
		return "", fmt.Errorf("could not open database: %s", err.Error())
	}
	defer r.Close()
	var result string
	for r.Next() {
		err = r.Scan(&result)
		if err != nil {
			return "", err
		}
	}
	return result, nil
}

func (d *DBClient) SaveTuitsToDB(tl *[]lks.TuitLike, ctx context.Context) (string, error) {
	lenTL := len(*tl)
	lastTL := ""
	for i := range *tl {
		t := (*tl)[lenTL-1-i]
		_, err := d.BunDB.NewInsert().
			Model(&t).Returning("id").Exec(ctx)
		if err != nil {
			return lastTL, err
		}
		lastTL = t.ID

		_, err = d.BunDB.NewInsert().
			Model(&t.Author).Ignore().Exec(ctx)
		if err != nil {
			return lastTL, err
		}
		for _, m := range t.MediaData {
			_, err = d.BunDB.NewInsert().
				Model(&m).Exec(ctx)
			if err != nil {
				return lastTL, err
			}
		}
	}
	return lastTL, nil
}

func ReadLineFromFile(filename string) (*[]lks.TuitLike, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	tlList := []lks.TuitLike{}
	for scanner.Scan() {
		line := scanner.Bytes()
		tl, err := lks.LineByteToTuitLike(line)
		if err != nil {
			log.Println(err)
			continue
		}
		tlList = append(tlList, *tl)
	}

	return &tlList, nil
}
