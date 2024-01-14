package types

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"go.universe.tf/garden/gogarden/forms"
)

type QRState int

//go:generate stringer -type=QRState

const (
	Missing QRState = iota
	Applied
	Ignore
)

func (b *QRState) Scan(value any) error {
	if v, ok := value.(string); ok {
		switch v {
		case "wanted":
			*b = Missing
		case "applied":
			*b = Applied
		case "none":
			*b = Ignore
		default:
			return fmt.Errorf("no conversion from %q to QRState", v)
		}
		return nil
	}
	return fmt.Errorf("no conversion from %v (%T) to QRState", value, value)
}

func (b QRState) Value() (driver.Value, error) {
	switch b {
	case Missing:
		return "wanted", nil
	case Applied:
		return "applied", nil
	case Ignore:
		return "none", nil
	default:
		return nil, errors.New("no conversion")
	}
}

func (b *QRState) UnmarshalText(bs []byte) error {
	switch string(bs) {
	case "Missing":
		*b = Missing
	case "Applied":
		*b = Applied
	case "Ignore":
		*b = Ignore
	default:
		return errors.New("no conversion")
	}
	return nil
}

func (QRState) SelectOptions() []forms.SelectOption {
	return []forms.SelectOption{
		{
			Value: "Missing",
			Label: "Missing",
		},
		{
			Value: "Applied",
			Label: "Applied",
		},
		{
			Value: "Ignore",
			Label: "Ignore",
		},
	}
}
