package plantlifespan

import (
	"database/sql/driver"
	"errors"

	"go.universe.tf/garden/forms"
)

type PlantLifespan int

//go:generate stringer -type=PlantLifespan

const (
	Unknown PlantLifespan = iota
	Annual
	Perennial
)

func (b *PlantLifespan) Scan(value any) error {
	if value == nil {
		*b = Unknown
		return nil
	}
	if v, ok := value.(string); ok {
		switch v {
		case "Annual":
			*b = Annual
		case "Perennial":
			*b = Perennial
		default:
			return errors.New("no conversion")
		}
		return nil
	}
	return errors.New("no conversion")
}

func (b PlantLifespan) Value() (driver.Value, error) {
	switch b {
	case Unknown:
		return nil, nil
	case Annual:
		return "Annual", nil
	case Perennial:
		return "Perennial", nil
	default:
		return nil, errors.New("no conversion")
	}
}

func (b PlantLifespan) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *PlantLifespan) UnmarshalText(bs []byte) error {
	switch string(bs) {
	case "Unknown":
		*b = Unknown
	case "Annual":
		*b = Annual
	case "Perennial":
		*b = Perennial
	default:
		return errors.New("no conversion")
	}
	return nil
}

func (PlantLifespan) SelectOptions() []forms.SelectOption {
	return []forms.SelectOption{
		{
			Value: "Unknown",
			Label: "Unknown",
		},
		{
			Value: "Annual",
			Label: "Annual",
		},
		{
			Value: "Perennial",
			Label: "Perennial",
		},
	}
}
