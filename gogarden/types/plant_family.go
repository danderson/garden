package types

import (
	"database/sql/driver"
	"errors"
)

//go:generate stringer -type=PlantFamily -trimprefix=PlantFamily

type PlantFamily int

const (
	PlantFamilyUnknown PlantFamily = iota
	PlantFamilyAdoxaceae
	PlantFamilyAllium
	PlantFamilyAmaranthaceae
	PlantFamilyApiaceae
	PlantFamilyApocynaceae
	PlantFamilyAsparagaceae
	PlantFamilyAsteraceae
	PlantFamilyBoraginaceae
	PlantFamilyBrassicaceae
	PlantFamilyCampanulaceae
	PlantFamilyCaprifoliaceae
	PlantFamilyCaryophyllaceae
	PlantFamilyConvolvulaceae
	PlantFamilyCucurbitaceae
	PlantFamilyFabaceae
	PlantFamilyLamiaceae
	PlantFamilyLinaceae
	PlantFamilyMalvaceae
	PlantFamilyOnagraceae
	PlantFamilyPapaveraceae
	PlantFamilyPoaceae
	PlantFamilyPolygonaceae
	PlantFamilyRanunculaceae
	PlantFamilyRosaceae
	PlantFamilySolanaceae
	PlantFamilyTropaeolaceae
	PlantFamilyViolaceae
	PlantFamilyWildflower
)

var (
	stringToPlantFamily = map[string]PlantFamily{}
	plantFamilyStrings  []string
)

func init() {
	for i := PlantFamilyUnknown; i <= PlantFamilyWildflower; i++ {
		stringToPlantFamily[i.String()] = i
		plantFamilyStrings = append(plantFamilyStrings, i.String())
	}
}

func (f *PlantFamily) Scan(value any) error {
	if value == nil {
		*f = PlantFamilyUnknown
		return nil
	}
	if v, ok := value.(string); ok && stringToPlantFamily[v] != PlantFamilyUnknown {
		*f = stringToPlantFamily[v]
		return nil
	}
	return errors.New("no conversion")
}

func (f PlantFamily) Value() (driver.Value, error) {
	if f == PlantFamilyUnknown {
		return nil, nil
	}
	return f.String(), nil
}

func PlantFamilyStrings() []string {
	return plantFamilyStrings
}
