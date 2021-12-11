package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/marcy-t/Golang-Logging-Sample/pkg/interfaces/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/marcy-t/Golang-Logging-Sample/pkg/logger"
)

type SqlHandler struct {
	Conn *sql.DB
}

type dbSettings struct {
	Host     string
	User     string
	Password string
	Database string
}

func NewHandler() (h db.SqlHandler, err error) {
	conf := dbSettings{
		Host:     "localhost",
		Database: "sample_db",
		User:     "user",
		Password: "password",
	}
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", conf.User, conf.Password, conf.Host, conf.Database)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		err = logger.GetApplicationError(err).Init("11", "Database configuration invalid.")
		return nil, err
	}
	myq := &SqlHandler{Conn: db}

	for retryCount := 10; retryCount > 0; retryCount-- {
		err = myq.Conn.Ping()
		if err == nil {
			logger.Info("10", "Connect successfully to database.")
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		err = logger.GetApplicationError(err).Init("12", "Connection failed to database.")
		return nil, err
	}
	h = myq
	return
}

/*
 * Receiver method of Sql Handler struct
 * Database Handler.
 */
func (h *SqlHandler) Exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	res, err = h.Conn.ExecContext(ctx, query, args...)
	return
}

func (h *SqlHandler) Query(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	rows, err = h.Conn.QueryContext(ctx, query, args...)
	return
}

func (h *SqlHandler) QueryRow(ctx context.Context, query string, args ...interface{}) (row *sql.Row) {
	row = h.Conn.QueryRowContext(ctx, query, args...)
	return
}

func (h *SqlHandler) ExecWithTx(txFunc func(*sql.Tx) error) (err error) {
	tx, err := h.Conn.Begin()
	if err != nil {
		err = logger.GetApplicationError(err).Init("xx", "Failed to start transaction.")
		return
	}

	defer func() {
		if p := recover(); p != nil {
			logger.Error(
				logger.NewApplicationError(p).
					Init("xx", "An error has occured. Transaction is rolled back..."),
			)
			err = tx.Rollback()
			panic(p)
		} else if err != nil {
			logger.Error(
				logger.GetApplicationError(err).
					Init("xx", "An error has occured. Transaction is rolled back..."),
			)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = rollbackErr
			}
		} else {
			logger.Debug("XX", "Begin transaction commit...")
			err = tx.Commit()
		}
	}()

	err = txFunc(tx)
	return
}
