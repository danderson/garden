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
		4: dropPlantsExplicitNulls,
		5: dropPlantLocationsExplicitNulls,
		7: fixupLifespanValues,
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

func dropPlantsExplicitNulls(tx *sqlx.Tx) error {
	var version int
	if err := tx.Get(&version, "PRAGMA schema_version"); err != nil {
		return err
	}
	if _, err := tx.Exec("update plants set name_from_seed=0 where name_from_seed is null"); err != nil {
		return err
	}
	if _, err := tx.Exec("PRAGMA writable_schema=ON"); err != nil {
		return err
	}
	var schema string
	if err := tx.Get(&schema, "select sql from sqlite_schema where type='table' and name='plants'"); err != nil {
		return err
	}
	replacements := map[string]string{
		`"name" TEXT`:              `"name" TEXT NOT NULL`,
		`"name_from_seed" INTEGER`: `"name_from_seed" INTEGER NOT NULL`,
	}
	for from, to := range replacements {
		schema = strings.ReplaceAll(schema, from, to)
	}
	if _, err := tx.Exec("update sqlite_schema set sql=? where type='table' and name='plants'", schema); err != nil {
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

func dropPlantLocationsExplicitNulls(tx *sqlx.Tx) error {
	var version int
	if err := tx.Get(&version, "PRAGMA schema_version"); err != nil {
		return err
	}
	if _, err := tx.Exec("PRAGMA writable_schema=ON"); err != nil {
		return err
	}
	var schema string
	if err := tx.Get(&schema, "select sql from sqlite_schema where type='table' and name='plant_locations'"); err != nil {
		return err
	}
	replacements := map[string]string{
		`"location_id" INTEGER`: `"location_id" INTEGER NOT NULL`,
		`"plant_id" INTEGER`:    `"plant_id" INTEGER NOT NULL`,
	}
	for from, to := range replacements {
		schema = strings.ReplaceAll(schema, from, to)
	}
	if _, err := tx.Exec("update sqlite_schema set sql=? where type='table' and name='plant_locations'", schema); err != nil {
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

func fixupLifespanValues(tx *sqlx.Tx) error {
	if _, err := tx.Exec("update seeds set lifespan=NULL where lifespan='U'"); err != nil {
		return err
	}
	if _, err := tx.Exec("update seeds set lifespan='Annual' where lifespan='A'"); err != nil {
		return err
	}
	if _, err := tx.Exec("update seeds set lifespan='Perennial' where lifespan='P'"); err != nil {
		return err
	}
	if _, err := tx.Exec("update seeds set lifespan='Perennial' where lifespan='B'"); err != nil {
		return err
	}
	return nil
}
