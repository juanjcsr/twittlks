package db

import (
	"bufio"
	"context"
	"database/sql"

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

func OpenSQLConn() *DBClient {
	dsn := "postgres://juan.sanchez:@localhost:5432/twitlks?sslmode=disable"

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

func (db *DBClient) CreateTables(ctx context.Context) {

	err := db.BunDB.ResetModel(ctx, (*lks.Users)(nil),
		(*lks.TuitLike)(nil), (*lks.Media)(nil))
	_, err = db.BunDB.NewCreateTable().
		Model((*lks.Users)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.BunDB.NewCreateTable().
		Model((*lks.TuitLike)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.BunDB.NewCreateTable().
		Model((*lks.Media)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}
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
