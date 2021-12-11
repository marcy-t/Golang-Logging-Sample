package database

import (
	"context"
	"database/sql"
)

type SqlHandler interface {
	Exec(context.Context, string, ...interface{}) (sql.Result, error)
	Query(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRow(context.Context, string, ...interface{}) *sql.Row
	ExecWithTx(txFunc func(*sql.Tx) error) error
}
