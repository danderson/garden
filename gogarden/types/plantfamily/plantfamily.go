package plantfamily

import (
	"database/sql/driver"
	"errors"
)

//go:generate stringer -type=PlantFamily

type PlantFamily int

const (
	Unknown PlantFamily = iota
	Adoxaceae
	Allium
	Amaranthaceae
	Apiaceae
	Apocynaceae
	Asparagaceae
	Asteraceae
	Boraginaceae
	Brassicaceae
	Campanulaceae
	Caprifoliaceae
	Caryophyllaceae
	Convolvulaceae
	Cucurbitaceae
	Fabaceae
	Lamiaceae
	Linaceae
	Malvaceae
	Onagraceae
	Papaveraceae
	Poaceae
	Polygonaceae
	Ranunculaceae
	Rosaceae
	Solanaceae
	Tropaeolaceae
	Violaceae
	Wildflower
)

var (
	stringToPlantFamily = map[string]PlantFamily{}
	plantFamilyStrings  []string
)

func init() {
	for i := Unknown; i <= Wildflower; i++ {
		stringToPlantFamily[i.String()] = i
		plantFamilyStrings = append(plantFamilyStrings, i.String())
	}
}

func (f *PlantFamily) Scan(value any) error {
	if value == nil {
		*f = Unknown
		return nil
	}
	if v, ok := value.(string); ok && stringToPlantFamily[v] != Unknown {
		*f = stringToPlantFamily[v]
		return nil
	}
	return errors.New("no conversion")
}

func (f PlantFamily) Value() (driver.Value, error) {
	if f == Unknown {
		return nil, nil
	}
	return f.String(), nil
}

func (f PlantFamily) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *PlantFamily) UnmarshalText(bs []byte) error {
	if v, ok := stringToPlantFamily[string(bs)]; ok {
		*f = v
		return nil
	}
	return errors.New("no conversion")
}

func (PlantFamily) SelectOptions() []string {
	return plantFamilyStrings
}

func Strings() []string {
	return plantFamilyStrings
}
