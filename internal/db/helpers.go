package db

import (
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	migrations "github.com/skynexus/soundx/db/sql"
	"github.com/skynexus/soundx/pkg/dbx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func pgStatementBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func idScan(row sq.RowScanner) (*int64, error) {
	var id int64
	err := row.Scan(&id)
	return &id, err
}

func rowsScan[T any](rows *sql.Rows, f func(s sq.RowScanner) (*T, error)) ([]T, error) {
	var values []T

	for rows.Next() {
		if value, err := f(rows); err == sql.ErrNoRows {
			return nil, nil
		} else if err != nil {
			return nil, err
		} else {
			values = append(values, *value)
		}
	}

	return values, nil
}

type Repository struct {
	dbx.Repository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{*dbx.NewRepository(db)}
}

type Resources struct {
	suite.Suite
	DB *sql.DB
}

func (r *Resources) SetupSuite() {
	dsn := dbx.ContainerTestInstance(r.T(), func(conf *dbx.PostgreSQLContainerConfig) {
		conf.User = "postgres"
		conf.Password = "tsetxdnuos"
		conf.Database = "soundxtest"
		conf.Migrations = migrations.Archive
	})

	db, openErr := sql.Open("postgres", dsn)
	require.NoError(r.T(), openErr, dsn)
	r.DB = db
}

func (r *Resources) TearDownSuite() {
	r.DB.Close()
}

func (r *Resources) BeforeTest(suiteName, testName string) {
	tables := []string{
		"sounds",
		"playlists",
	}

	query := fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))

	_, err := r.DB.Exec(query)
	require.NoError(r.T(), err)
}
