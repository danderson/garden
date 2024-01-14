package tribool

import (
	"database/sql/driver"
	"errors"
)

type Tribool int

//go:generate stringer -type=Tribool

const (
	Unknown Tribool = iota
	True
	False
)

func (b *Tribool) Scan(value any) error {
	if value == nil {
		*b = Unknown
		return nil
	}
	if v, ok := value.(int64); ok {
		if v == 0 {
			*b = False
		} else {
			*b = True
		}
		return nil
	}
	return errors.New("no conversion")
}

func (b Tribool) Value() (driver.Value, error) {
	switch b {
	case Unknown:
		return nil, nil
	case True:
		return int64(1), nil
	case False:
		return int64(0), nil
	default:
		return nil, errors.New("no conversion")
	}
}

func (b Tribool) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *Tribool) UnmarshalText(bs []byte) error {
	switch string(bs) {
	case "Unknown":
		*b = Unknown
	case "True":
		*b = True
	case "False":
		*b = False
	default:
		return errors.New("no conversion")
	}
	return nil
}

func (Tribool) SelectOptions() []string {
	return []string{"Unknown", "True", "False"}
}
