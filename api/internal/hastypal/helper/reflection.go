package helper

import (
	"reflect"
)

func HasField[T interface{}](obj T, name string) bool {
	structType := reflect.TypeOf(obj)

	structVal := reflect.ValueOf(obj)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		if fieldName == name && field.IsZero() {
			return false
		}

		return true
	}

	return false
}
