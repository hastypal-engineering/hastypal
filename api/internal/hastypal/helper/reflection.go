package helper

import (
	"reflect"
)

type ReflectionHelper[T interface{}] struct{}

func NewReflectionHelper[T interface{}]() *ReflectionHelper[T] {
	return &ReflectionHelper[T]{}
}

func (*ReflectionHelper[T]) HasField(obj T, name string) bool {
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

func (*ReflectionHelper[T]) Merge(actual T, updated T) T {
	merged := actual

	structType := reflect.TypeOf(actual)

	structVal := reflect.ValueOf(actual)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		newValue := reflect.ValueOf(updated).FieldByName(fieldName)

		if !field.IsZero() {
			reflect.ValueOf(&merged).Elem().FieldByName(fieldName).Set(newValue)
		}
	}

	return merged
}
