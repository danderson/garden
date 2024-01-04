package types

import (
	"database/sql/driver"
	"errors"
	"time"
)

type Time struct {
	time.Time
}

func (n *Time) Scan(value any) error {
	if value == nil {
		*n = Time{}
		return nil
	}
	if v, ok := value.(int64); ok {
		*n = Time{time.Unix(0, v).UTC()}
		return nil
	}
	return errors.New("no conversion")
}

func (n *Time) Value() (driver.Value, error) {
	if n.IsZero() {
		return int64(0), nil
	}
	return n.UnixNano(), nil
}
