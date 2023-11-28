package dbx

import (
	"context"
	"database/sql"
	"errors"
)

var ErrUnsupportedDatabaseHandle = errors.New("unsupported database handle for transactions")

type Repository struct {
	handle *sql.DB
	dbtx   DBTx
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{handle: db, dbtx: db}
}

func (repo *Repository) Clone() *Repository {
	return &Repository{handle: repo.handle, dbtx: repo.handle}
}

func (repo *Repository) Db() DBTx {
	return repo.dbtx
}

func (repo *Repository) Transactional(callback func() error) error {
	switch handle := repo.dbtx.(type) {
	case *sql.Tx:
		return callback()
	case *sql.DB:
		ctx := context.Background()
		tx, txErr := handle.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  false,
		})
		if txErr != nil {
			return txErr
		}

		repo.dbtx = tx
		defer func() {
			repo.dbtx = repo.handle
			tx.Rollback()
		}()

		if callbackErr := callback(); callbackErr != nil {
			return callbackErr
		}

		if commitErr := tx.Commit(); commitErr != nil {
			return commitErr
		}

		return nil
	default:
		return ErrUnsupportedDatabaseHandle
	}
}
