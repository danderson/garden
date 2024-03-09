package soiltype

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"go.universe.tf/garden/forms"
)

type SoilType int

//go:generate stringer -type=SoilType

const (
	Unknown SoilType = iota
	Dry
	Wet
	Both
)

func (b *SoilType) Scan(value any) error {
	if v, ok := value.(string); ok {
		switch v {
		case "Unknown":
			*b = Unknown
		case "Dry":
			*b = Dry
		case "Wet":
			*b = Wet
		case "Both":
			*b = Both
		default:
			return fmt.Errorf("no conversion from %q to SoilType", v)
		}
		return nil
	}
	return fmt.Errorf("no conversion from %v (%T) to SoilType", value, value)
}

func (b SoilType) Value() (driver.Value, error) {
	switch b {
	case Unknown, Dry, Wet, Both:
		return b.String(), nil
	default:
		return nil, errors.New("no conversion")
	}
}

func (b *SoilType) UnmarshalText(bs []byte) error {
	switch string(bs) {
	case "Unknown":
		*b = Unknown
	case "Dry":
		*b = Dry
	case "Wet":
		*b = Wet
	case "Both":
		*b = Both
	default:
		return errors.New("no conversion")
	}
	return nil
}

func (b SoilType) MarshalText() ([]byte, error) {
	switch b {
	case Unknown, Dry, Wet, Both:
		return []byte(b.String()), nil
	default:
		return nil, errors.New("no conversion")
	}
}

func (SoilType) SelectOptions() []forms.SelectOption {
	return []forms.SelectOption{
		{
			Value: "Unknown",
			Label: "Unknown",
		},
		{
			Value: "Dry",
			Label: "Dry",
		},
		{
			Value: "Wet",
			Label: "Wet",
		},
		{
			Value: "Both",
			Label: "Both",
		},
	}
}
