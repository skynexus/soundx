// Package sql contains all migrations embedded in this project.
package sql

import "embed"

//go:embed *.sql
var Archive embed.FS
