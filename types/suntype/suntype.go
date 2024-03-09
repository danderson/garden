package suntype

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"go.universe.tf/garden/forms"
)

type SunType int

//go:generate stringer -type=SunType

const (
	Unknown SunType = iota
	Full
	Partial
	Shade
)

func (b *SunType) Scan(value any) error {
	if v, ok := value.(string); ok {
		switch v {
		case "Unknown":
			*b = Unknown
		case "Full":
			*b = Full
		case "Partial":
			*b = Partial
		case "Shade":
			*b = Shade
		default:
			return fmt.Errorf("no conversion from %q to SunType", v)
		}
		return nil
	}
	return fmt.Errorf("no conversion from %v (%T) to SunType", value, value)
}

func (b SunType) Value() (driver.Value, error) {
	switch b {
	case Unknown, Full, Partial, Shade:
		return b.String(), nil
	default:
		return nil, errors.New("no conversion")
	}
}

func (b *SunType) UnmarshalText(bs []byte) error {
	switch string(bs) {
	case "Unknown":
		*b = Unknown
	case "Full":
		*b = Full
	case "Partial":
		*b = Partial
	case "Shade":
		*b = Shade
	default:
		return errors.New("no conversion")
	}
	return nil
}

func (b SunType) MarshalText() ([]byte, error) {
	switch b {
	case Unknown, Full, Partial, Shade:
		return []byte(b.String()), nil
	default:
		return nil, errors.New("no conversion")
	}
}

func (SunType) SelectOptions() []forms.SelectOption {
	return []forms.SelectOption{
		{
			Value: "Unknown",
			Label: "Unknown",
		},
		{
			Value: "Full",
			Label: "Full",
		},
		{
			Value: "Partial",
			Label: "Partial",
		},
		{
			Value: "Shade",
			Label: "Shade",
		},
	}
}
