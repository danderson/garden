package plantfamily

import (
	"database/sql/driver"
	"errors"

	"go.universe.tf/garden/forms"
)

//go:generate stringer -type=PlantFamily

type PlantFamily int

const (
	Unknown PlantFamily = iota
	Adoxaceae
	Aizoaceae
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
	Lilliaceae
	Linaceae
	Malvaceae
	Onagraceae
	Papaveraceae
	Plantaginaceae
	Plumbaginaceae
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
	plantFamilyOptions  []forms.SelectOption
)

func init() {
	for i := Unknown; i <= Wildflower; i++ {
		stringToPlantFamily[i.String()] = i
		plantFamilyOptions = append(plantFamilyOptions, forms.SelectOption{
			Value: i.String(),
			Label: i.String(),
		})
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

func (PlantFamily) SelectOptions() []forms.SelectOption {
	return plantFamilyOptions
}
