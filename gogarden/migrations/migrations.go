package migrations

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/jmoiron/sqlx"
)

//go:embed *.sql
var migrations embed.FS

func FileMigrations() fs.FS {
	return migrations
}

func GoMigrations() map[int]func(*sqlx.Tx) error {
	return map[int]func(*sqlx.Tx) error{
		2: dropSeedExplicitNulls,
		3: dropLocationExplicitNulls,
	}
}

func dropSeedExplicitNulls(tx *sqlx.Tx) error {
	var version int
	if err := tx.Get(&version, "PRAGMA schema_version"); err != nil {
		return err
	}
	if _, err := tx.Exec("PRAGMA writable_schema=ON"); err != nil {
		return err
	}
	var schema string
	if err := tx.Get(&schema, "select sql from sqlite_schema where type='table' and name='seeds'"); err != nil {
		return err
	}
	replacements := map[string]string{
		"INTEGER NULL": "INTEGER",
		`"name" TEXT`:  `"name" TEXT NOT NULL`,
	}
	for from, to := range replacements {
		schema = strings.ReplaceAll(schema, from, to)
	}
	if _, err := tx.Exec("update sqlite_schema set sql=? where type='table' and name='seeds'", schema); err != nil {
		return err
	}
	if _, err := tx.Exec(fmt.Sprintf("PRAGMA schema_version=%d", version+1)); err != nil {
		return err
	}
	if _, err := tx.Exec("PRAGMA writable_schema=OFF"); err != nil {
		return err
	}
	if _, err := tx.Exec("PRAGMA integrity_check"); err != nil {
		return err
	}
	return nil
}

func dropLocationExplicitNulls(tx *sqlx.Tx) error {
	var version int
	if err := tx.Get(&version, "PRAGMA schema_version"); err != nil {
		return err
	}
	if _, err := tx.Exec("PRAGMA writable_schema=ON"); err != nil {
		return err
	}
	var schema string
	if err := tx.Get(&schema, "select sql from sqlite_schema where type='table' and name='locations'"); err != nil {
		return err
	}
	replacements := map[string]string{
		`"name" TEXT`: `"name" TEXT NOT NULL`,
	}
	for from, to := range replacements {
		schema = strings.ReplaceAll(schema, from, to)
	}
	if _, err := tx.Exec("update sqlite_schema set sql=? where type='table' and name='locations'", schema); err != nil {
		return err
	}
	if _, err := tx.Exec(fmt.Sprintf("PRAGMA schema_version=%d", version+1)); err != nil {
		return err
	}
	if _, err := tx.Exec("PRAGMA writable_schema=OFF"); err != nil {
		return err
	}
	if _, err := tx.Exec("PRAGMA integrity_check"); err != nil {
		return err
	}
	return nil
}
