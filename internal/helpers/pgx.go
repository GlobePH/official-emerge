package helpers

import (
	"github.com/jackc/pgx"
)

type AfterConnectFunc func(*pgx.Conn) error

func NewPgxPool(uri string, afterConnect AfterConnectFunc) (pool *pgx.ConnPool, err error) {
	cfg, err := pgx.ParseURI(uri)
	if err != nil {
		return
	}
	return pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   cfg,
		AfterConnect: afterConnect,
	})
}
