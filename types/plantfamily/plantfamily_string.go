// Code generated by "stringer -type=PlantFamily"; DO NOT EDIT.

package plantfamily

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[Adoxaceae-1]
	_ = x[Aizoaceae-2]
	_ = x[Allium-3]
	_ = x[Amaranthaceae-4]
	_ = x[Apiaceae-5]
	_ = x[Apocynaceae-6]
	_ = x[Asparagaceae-7]
	_ = x[Asteraceae-8]
	_ = x[Boraginaceae-9]
	_ = x[Brassicaceae-10]
	_ = x[Campanulaceae-11]
	_ = x[Caprifoliaceae-12]
	_ = x[Caryophyllaceae-13]
	_ = x[Convolvulaceae-14]
	_ = x[Cucurbitaceae-15]
	_ = x[Fabaceae-16]
	_ = x[Lamiaceae-17]
	_ = x[Lilliaceae-18]
	_ = x[Linaceae-19]
	_ = x[Malvaceae-20]
	_ = x[Onagraceae-21]
	_ = x[Papaveraceae-22]
	_ = x[Plantaginaceae-23]
	_ = x[Plumbaginaceae-24]
	_ = x[Poaceae-25]
	_ = x[Polygonaceae-26]
	_ = x[Ranunculaceae-27]
	_ = x[Rosaceae-28]
	_ = x[Solanaceae-29]
	_ = x[Tropaeolaceae-30]
	_ = x[Violaceae-31]
	_ = x[Wildflower-32]
}

const _PlantFamily_name = "UnknownAdoxaceaeAizoaceaeAlliumAmaranthaceaeApiaceaeApocynaceaeAsparagaceaeAsteraceaeBoraginaceaeBrassicaceaeCampanulaceaeCaprifoliaceaeCaryophyllaceaeConvolvulaceaeCucurbitaceaeFabaceaeLamiaceaeLilliaceaeLinaceaeMalvaceaeOnagraceaePapaveraceaePlantaginaceaePlumbaginaceaePoaceaePolygonaceaeRanunculaceaeRosaceaeSolanaceaeTropaeolaceaeViolaceaeWildflower"

var _PlantFamily_index = [...]uint16{0, 7, 16, 25, 31, 44, 52, 63, 75, 85, 97, 109, 122, 136, 151, 165, 178, 186, 195, 205, 213, 222, 232, 244, 258, 272, 279, 291, 304, 312, 322, 335, 344, 354}

func (i PlantFamily) String() string {
	if i < 0 || i >= PlantFamily(len(_PlantFamily_index)-1) {
		return "PlantFamily(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PlantFamily_name[_PlantFamily_index[i]:_PlantFamily_index[i+1]]
}
