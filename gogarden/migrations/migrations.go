package migrations

import (
	"embed"
	"io/fs"

	"github.com/jmoiron/sqlx"
)

//go:embed *.sql
var migrations embed.FS

func FileMigrations() fs.FS {
	return migrations
}

func GoMigrations() map[int]func(*sqlx.Tx) error {
	return map[int]func(*sqlx.Tx) error{}
}
