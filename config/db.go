package config

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	RequestTable  = "requests"
	ResponseTable = "responses"
	TxHash        = "tx_hash"
	User          = "user"
	RequestId     = "request_id"
	BlockNumber   = "block_number"
)

var (
	db  *bun.DB
	ctx = context.Background()
)

type Request struct {
	Id          int    `bun:"id,autoincrement"`
	BlockNumber int64  `bun:"block_number,notnull"`
	TxHash      string `bun:"tx_hash,pk,notnull"`
	TxIndex     int    `bun:"tx_index,notnull"`
	Amount      int    `bun:"amount,notnull"`
	User        string `bun:"user,notnull"`
	RequestId   string `bun:"request_id,notnull"`
}

type Response struct {
	Id          int    `bun:"id,autoincrement"`
	BlockNumber int64  `bun:"block_number,notnull"`
	TxHash      string `bun:"tx_hash,pk,notnull"`
	TxIndex     int    `bun:"tx_index,notnull"`
	User        string `bun:"user,notnull"`
	RequestId   string `bun:"request_id,notnull"`
	PrizeIds    []int  `bun:"prize_ids,notnull" json:"prizeIds"`
}

func ConnectDb(DSN string) {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(DSN)))
	db = bun.NewDB(sqlDb, pgdialect.New())
}

func CreateTable() error {
	_, err := db.NewCreateTable().Model((*Request)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*Response)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func CreateIndexes() error {
	_, err := db.NewCreateIndex().Table(RequestTable).Column(TxHash).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateIndex().Table(RequestTable).Column(User).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateIndex().Table(RequestTable).Column(RequestId).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateIndex().Table(RequestTable).Column(BlockNumber).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateIndex().Table(ResponseTable).Column(TxHash).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateIndex().Table(ResponseTable).Column(User).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateIndex().Table(ResponseTable).Column(RequestId).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
