package config

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	Table       = "transactions"
	From        = "from"
	To          = "to"
	BlockNumber = "block_number"
)

var (
	db  *bun.DB
	ctx = context.Background()
)

type Transaction struct {
	Id           int     `bun:"id,pk,autoincrement"`
	BlockNumber  int64   `bun:"block_number,notnull"`
	Balance      float64 `bun:"balance,notnull"`
	RawBalance   string  `bun:"raw_balance,notnull"`
	Hash         string  `bun:"hash,notnull"`
	Amount       float64 `bun:"amount,notnull"`
	RawAmount    string  `bun:"raw_amount,notnull"`
	From         string  `bun:"from,notnull"`
	To           string  `bun:"to"`
	TokenAddress string  `bun:"token_address,notnull"`
	Token        string  `bun:"token,notnull"`
}

func ConnectDb(DSN string) {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(DSN)))
	db = bun.NewDB(sqlDb, pgdialect.New())
}

func CreateTable() error {
	_, err := db.NewCreateTable().Model((*Transaction)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func CreateIndexes() error {
	_, err := db.NewCreateIndex().Table(Table).Column(From).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateIndex().Table(Table).Column(To).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateIndex().Table(Table).Column(BlockNumber).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
