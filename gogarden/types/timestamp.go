package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
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

type TextTime struct {
	time.Time
}

func (n *TextTime) Scan(value any) error {
	if value == nil {
		*n = TextTime{}
		return nil
	}
	if v, ok := value.(string); ok {
		if t, err := time.Parse("2006-01-02T15:04:05", v); err == nil {
			*n = TextTime{t}
			return nil
		}
		if t, err := time.Parse("2006-01-02T15:04:05.999999", v); err == nil {
			*n = TextTime{t}
			return nil
		}
		return fmt.Errorf("Unparseable time %q", v)
	}

	return errors.New("no conversion")
}

func (n *TextTime) Value() (driver.Value, error) {
	if n.IsZero() {
		return "", nil
	}
	return n.Time.Format("2006-01-02T15:04:05.999999"), nil
}
