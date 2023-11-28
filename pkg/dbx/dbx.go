// Package dbx implements functionality that help working with databases.
package dbx

import (
	"context"
	"database/sql"
	"log"
	"net/url"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	AWS_POSTGRES_TIMEOUT = 2 * time.Second
)

// https://github.com/golang/go/issues/14468
type DBTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func Open(ctx context.Context, dsn *url.URL) (*sql.DB, error) {
	addTimeoutToDSN(dsn, AWS_POSTGRES_TIMEOUT)
	log.Printf("dsn: %s", dsn.Redacted())

	db, openErr := sql.Open("postgres", dsn.String())
	if openErr != nil {
		return nil, openErr
	}

	// NOTE: this has no practical effect at the moment due to a bug in
	// the postgres driver: https://github.com/lib/pq/issues/1020
	//
	// We work around it by adding the connect_timeout parameter to the
	// connection string (see addTimeoutToDSN for more information).
	// Leave this code here - ctx will kick in once the bug is fixed.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if pingErr := db.PingContext(ctx); pingErr != nil {
		return nil, pingErr
	} else {
		return db, nil
	}
}

// addTimeoutToDSN adds the connect_timeout parameter to the connection
// string. This is used as a workaround for a bug in the postgres driver
// that prevents timeout context from working correctly.
//
// See also: https://github.com/lib/pq/issues/1020
func addTimeoutToDSN(dsn *url.URL, timeout time.Duration) {
	secs := int64(timeout / time.Second)
	q := dsn.Query()
	q.Set("connect_timeout", strconv.FormatInt(secs, 10))
	dsn.RawQuery = q.Encode()
}
