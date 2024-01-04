package migrations

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.universe.tf/garden/gogarden/types"
)

func moveBoxInventory(tx *sqlx.Tx) error {
	type box struct {
		ID        int
		Name      string
		WantQR    bool `db:"want_qr"`
		QRApplied bool `db:"qr_applied"`
	}
	boxes := []box{}
	if err := tx.Select(&boxes, "select * from boxinventory_box"); err != nil {
		return fmt.Errorf("getting boxes: %w", err)
	}

	for _, box := range boxes {
		qrState := types.QRStateIgnore
		if box.QRApplied {
			qrState = types.QRStateApplied
		} else if box.WantQR {
			qrState = types.QRStateWant
		}
		if _, err := tx.Exec("insert into location (id, name, qr_state) values (?,?,?)", box.ID, box.Name, qrState); err != nil {
			return fmt.Errorf("inserting box ID %d: %w", box.ID, err)
		}
	}

	type content struct {
		ID      int
		Name    string
		BoxID   int `db:"box_id"`
		Planted time.Time
		Removed sql.NullTime
	}
	contents := []content{}
	if err := tx.Select(&contents, "select * from boxinventory_boxcontent"); err != nil {
		return fmt.Errorf("getting box contents: %w", err)
	}

	for _, plant := range contents {
		removed := time.Time{}
		if plant.Removed.Valid {
			removed = plant.Removed.Time
		}
		if _, err := tx.Exec("insert into planted (id, name, location, planted, removed) values (?,?,?,?,?)", plant.ID, plant.Name, plant.BoxID, plant.Planted.UTC().UnixNano(), removed.UTC().UnixNano()); err != nil {
			return fmt.Errorf("inserting plant %d: %w", plant.ID, err)
		}
	}

	return nil
}
