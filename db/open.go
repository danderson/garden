package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func Open(path string, fileMigrations fs.FS, goMigrations map[int]func(*sqlx.Tx) error) (*DB, error) {
	if path == "" {
		path = ":memory:"
	}
	url := fmt.Sprintf("file:%s?_foreign_keys=true&_journal_mode=WAL&loc=UTC", path)
	db, err := sqlx.Open("sqlite", url)
	if err != nil {
		return nil, err
	}

	var dbVersion int
	if err := db.Get(&dbVersion, "PRAGMA user_version"); err != nil {
		db.Close()
		return nil, fmt.Errorf("getting migration version: %w", err)
	}

	dbVersion, migrations, err := assembleMigrations(dbVersion, fileMigrations, goMigrations)
	if err != nil {
		return nil, err
	}

	log.Printf("database schema is version %d, have %d migrations to run", dbVersion, len(migrations))

	if len(migrations) > 0 {
		for _, m := range migrations {
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
			dbVersion++
			if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version=%d", dbVersion)); err != nil {
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

func assembleMigrations(dbVersion int, fileMigrations fs.FS, goMigrations map[int]func(*sqlx.Tx) error) (baseVersion int, fns []func(*sqlx.Tx) error, err error) {
	migrations := map[int]func(*sqlx.Tx) error{}
	lowest := math.MaxInt

	ents, err := fs.ReadDir(fileMigrations, ".")
	if err != nil {
		return 0, nil, fmt.Errorf("listing migration files: %w", err)
	}
	for _, ent := range ents {
		fields := strings.Split(ent.Name(), "_")
		if len(fields) < 1 {
			return 0, nil, fmt.Errorf("unparseable filename %q", ent.Name())
		}
		i, err := strconv.Atoi(fields[0])
		if err != nil {
			return 0, nil, fmt.Errorf("filename %q doesn't begin with a number", ent.Name())
		}
		if i < 0 {
			return 0, nil, fmt.Errorf("filename %q has invalid migration number %d, must be 1 or more", ent.Name(), i)
		}
		if i <= dbVersion {
			continue
		}
		if migrations[i] != nil {
			return 0, nil, fmt.Errorf("duplicate migration number %d", i)
		}
		bs, err := fs.ReadFile(fileMigrations, ent.Name())
		if err != nil {
			return 0, nil, fmt.Errorf("reading migration %q: %w", ent.Name(), err)
		}
		m := &sqlMigration{
			filename: ent.Name(),
			sql:      string(bs),
		}
		migrations[i] = m.do
		lowest = min(lowest, i)
	}
	for i, f := range goMigrations {
		if i <= dbVersion {
			continue
		}
		if migrations[i] != nil {
			return 0, nil, fmt.Errorf("duplicate go migration number %d", i)
		}
		migrations[i] = f
		lowest = min(lowest, i)
	}

	if len(migrations) == 0 {
		return dbVersion, nil, nil
	}

	if dbVersion == 0 {
		dbVersion = lowest - 1
	} else if lowest > dbVersion {
		return 0, nil, fmt.Errorf("database version %d is older than earliest known migration %d, cannot upgrade", dbVersion, lowest)
	}

	ret := make([]func(*sqlx.Tx) error, 0, len(migrations))
	for i := range len(migrations) {
		version := dbVersion + i + 1
		if migrations[version] == nil {
			return 0, nil, fmt.Errorf("missing migration number %d", version)
		}
		ret = append(ret, migrations[version])
	}
	return dbVersion, ret, nil
}
