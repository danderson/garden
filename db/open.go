package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func Open(path string, fileMigrations fs.FS, goMigrations map[int]func(*sqlx.Tx) error) (*DB, error) {
	migrations, err := assembleMigrations(fileMigrations, goMigrations)
	if err != nil {
		return nil, err
	}

	if path == "" {
		path = ":memory:"
	}
	url := fmt.Sprintf("file:%s?_foreign_keys=true&_journal_mode=WAL&loc=UTC", path)
	db, err := sqlx.Open("sqlite", url)
	if err != nil {
		return nil, err
	}

	var migration_version int
	if err := db.Get(&migration_version, "PRAGMA user_version"); err != nil {
		db.Close()
		return nil, fmt.Errorf("getting migration version: %w", err)
	}

	if migration_version > len(migrations) {
		db.Close()
		return nil, fmt.Errorf("database version %d exceeds max known version %d", migration_version, len(migrations))
	}

	if migration_version != len(migrations) {
		log.Printf("running migrations from version %d to %d", migration_version, len(migrations))
		for i, m := range migrations[migration_version:] {
			tx, err := db.Beginx()
			if err != nil {
				db.Close()
				return nil, fmt.Errorf("starting DB migration transaction: %w", err)
			}
			defer tx.Rollback()

			if err := m(tx); err != nil {
				db.Close()
				return nil, err
			}

			// Setting pragmas in sqlite with '?' substitutions apparently
			// doesn't work, so do scary sql-injecty formatting by hand.
			if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version=%d", migration_version+i+1)); err != nil {
				db.Close()
				return nil, fmt.Errorf("updating database migration version: %w", err)
			}
			if err := tx.Commit(); err != nil {
				db.Close()
				return nil, fmt.Errorf("committing schema migrations: %w", err)
			}
		}
	}

	return &DB{db, New(db)}, nil
}

type DB struct {
	*sqlx.DB
	*Queries
}

type Tx struct {
	*Queries
	*sqlx.Tx
}

func (db *DB) Tx(ctx context.Context) (*Tx, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{db.WithTx(tx.Tx), tx}, nil
}

func (db *DB) ReadTx(ctx context.Context) (*Tx, error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		return nil, err
	}
	return &Tx{db.WithTx(tx.Tx), tx}, nil
}

type sqlMigration struct {
	filename string
	sql      string
}

func (m *sqlMigration) do(tx *sqlx.Tx) error {
	if _, err := tx.Exec(m.sql); err != nil {
		return fmt.Errorf("executing migration %q: %w", m.filename, err)
	}
	return nil
}

func assembleMigrations(fileMigrations fs.FS, goMigrations map[int]func(*sqlx.Tx) error) ([]func(*sqlx.Tx) error, error) {
	migrations := map[int]func(*sqlx.Tx) error{}

	ents, err := fs.ReadDir(fileMigrations, ".")
	if err != nil {
		return nil, fmt.Errorf("listing migration files: %w", err)
	}
	for _, ent := range ents {
		fields := strings.Split(ent.Name(), "_")
		if len(fields) < 1 {
			return nil, fmt.Errorf("unparseable filename %q", ent.Name())
		}
		i, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("filename %q doesn't begin with a number", ent.Name())
		}
		if i < 0 {
			return nil, fmt.Errorf("filename %q has invalid migration number %d, must be 1 or more", ent.Name(), i)
		}
		if migrations[i] != nil {
			return nil, fmt.Errorf("duplicate migration number %d", i)
		}
		bs, err := fs.ReadFile(fileMigrations, ent.Name())
		if err != nil {
			return nil, fmt.Errorf("reading migration %q: %w", ent.Name(), err)
		}
		m := &sqlMigration{
			filename: ent.Name(),
			sql:      string(bs),
		}
		migrations[i] = m.do
	}
	for i, f := range goMigrations {
		if migrations[i] != nil {
			return nil, fmt.Errorf("duplicate go migration number %d", i)
		}
		migrations[i] = f
	}

	max := len(migrations)
	ret := make([]func(*sqlx.Tx) error, 0, max)
	for i := 1; i <= max; i++ {
		if migrations[i] == nil {
			return nil, fmt.Errorf("missing migration number %d", i)
		}
		ret = append(ret, migrations[i])
	}
	return ret, nil
}
