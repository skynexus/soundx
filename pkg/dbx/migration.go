package dbx

import (
	"database/sql"
	"errors"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

var (
	ErrNoChange   = errors.New("no change")
	ErrNilVersion = errors.New("no migration")
)

type logger struct{}

func (l *logger) Printf(format string, v ...interface{}) {
	format = "INFO: migration: " + format
	log.Printf(format, v...)
}

func (l *logger) Verbose() bool {
	return true
}

func migrator(db *sql.DB, fsys fs.FS) (*migrate.Migrate, error) {
	sourceDriver, fsErr := iofs.New(fsys, ".")
	if fsErr != nil {
		return nil, fsErr
	}

	driver, dbErr := postgres.WithInstance(db, &postgres.Config{})
	if dbErr != nil {
		return nil, dbErr
	}

	return migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
}

func Migrate(db *sql.DB, fsys fs.FS) error {
	m, mErr := migrator(db, fsys)
	if mErr != nil {
		return mErr
	}

	m.Log = &logger{}

	if upErr := m.Up(); upErr != nil {
		if errors.Is(upErr, migrate.ErrNoChange) {
			return ErrNoChange
		} else if errors.Is(upErr, migrate.ErrNilVersion) {
			return ErrNilVersion
		} else {
			return upErr
		}
	}

	return nil
}

type emptyFS struct{}

func (eFS emptyFS) Open(name string) (fs.File, error) {
	return nil, errors.New("not implemented")
}

func (eFS emptyFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return nil, nil
}

func Version(db *sql.DB) (version uint, dirty bool, err error) {
	m, mErr := migrator(db, emptyFS{})
	if mErr != nil {
		return version, dirty, mErr
	}

	m.Log = &logger{}

	return m.Version()
}
