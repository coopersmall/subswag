package postgres

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresClient(
	url string,
	opts *struct {
		MaxOpenConns int
		MaxIdleConns int
	},
) (*sql.DB, error) {
	config, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	conn := stdlib.OpenDB(*config)
	if opts != nil {
		conn.SetMaxOpenConns(opts.MaxOpenConns)
		conn.SetMaxIdleConns(opts.MaxIdleConns)
	}
	return conn, nil
}
