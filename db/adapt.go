package db

import "reflect"

func Adapt[To any, From any](from From) To {
	var ret To
	fromv := reflect.ValueOf(from)
	tov := reflect.ValueOf(&ret).Elem()
	t := tov.Type()
	for i := 0; i < t.NumField(); i++ {
		fi := t.Field(i)
		of := fromv.FieldByName(fi.Name)
		if !of.IsValid() {
			continue
		}
		tov.Field(i).Set(of)
	}
	return ret
}
