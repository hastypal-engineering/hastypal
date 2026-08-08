package reflection

import (
	"reflect"
)

func HasField(obj interface{}, name string) bool {
	hasField := false

	structType := reflect.TypeOf(obj)

	structVal := reflect.ValueOf(obj)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		if fieldName == name && field.IsZero() {
			hasField = false

			break
		}

		hasField = true
	}

	return hasField
}

func Merge(actual interface{}, updated interface{}) interface{} {
	merged := actual

	structType := reflect.TypeOf(actual)

	actualStructVal := reflect.ValueOf(actual)
	updatedStructVal := reflect.ValueOf(updated)
	fieldNum := actualStructVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := actualStructVal.Field(i)
		fieldName := structType.Field(i).Name

		newValue := reflect.ValueOf(updated).FieldByName(fieldName)

		if field.IsZero() && !updatedStructVal.Field(i).IsZero() {
			reflect.ValueOf(&merged).Elem().FieldByName(fieldName).Set(newValue)
		}
	}

	return merged
}
