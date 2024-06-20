package data

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"user/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewDatabase)

// Data .
type Data struct {
	db *pgxpool.Pool
}

// NewData .
func NewData(
	c *conf.Data,
	logger log.Logger,
	db *pgxpool.Pool,
	) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db: db,
	}, cleanup, nil
}

func NewDatabase(c *conf.Data) *pgxpool.Pool {
	conn, err := pgxpool.New(context.Background(), c.Database.Source)
	if err != nil {
		panic(err)
	}
	return conn
}
