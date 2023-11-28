package soundx

import (
	"database/sql"

	"github.com/skynexus/soundx/internal/db"
)

type server struct {
	r *db.Repository
}

func newServer(h *sql.DB) *server {
	return &server{
		r: db.NewRepository(h),
	}
}
