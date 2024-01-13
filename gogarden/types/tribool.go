package types

import (
	"database/sql/driver"
	"errors"
)

type Tribool int

//go:generate stringer -type=Tribool -trimprefix=Tribool

const (
	TriboolUnknown Tribool = iota
	TriboolTrue
	TriboolFalse
)

func (b *Tribool) Scan(value any) error {
	if value == nil {
		*b = TriboolUnknown
		return nil
	}
	if v, ok := value.(int64); ok {
		if v == 0 {
			*b = TriboolFalse
		} else {
			*b = TriboolTrue
		}
		return nil
	}
	return errors.New("no conversion")
}

func (b Tribool) Value() (driver.Value, error) {
	switch b {
	case TriboolUnknown:
		return nil, nil
	case TriboolTrue:
		return int64(1), nil
	case TriboolFalse:
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
		*b = TriboolUnknown
	case "True":
		*b = TriboolTrue
	case "False":
		*b = TriboolFalse
	default:
		return errors.New("no conversion")
	}
	return nil
}
