package helper

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"reflect"
)

type ReflectionHelper struct{}

func NewReflectionHelper() *ReflectionHelper {
	return &ReflectionHelper{}
}

func (*ReflectionHelper) HasField(obj interface{}, name string) bool {
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

func (*ReflectionHelper) Merge(actual interface{}, updated interface{}) interface{} {
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

func (*ReflectionHelper) ExtractDatabaseFields(entity interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(entity)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, exception.
			New("Entity provided is not a struct pointer").
			Trace("reflect.ValueOf", "reflection.go")
	}

	structValue := v.Elem()
	structType := structValue.Type()

	scanTargets := make([]interface{}, structType.NumField())

	fieldMap := make(map[string]int)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		dbTag := field.Tag.Get("db")

		if dbTag == "" {
			continue
		}

		fieldMap[dbTag] = i

		scanTargets[i] = structValue.Field(i).Addr().Interface()
	}

	return scanTargets, nil
}

func (*ReflectionHelper) ExtractTableName(entity interface{}) (string, error) {
	v := reflect.ValueOf(entity)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", exception.
			New("Entity provided is not a struct pointer").
			Trace("reflect.ValueOf", "reflection.go")
	}

	structValue := v.Elem()
	structType := structValue.Type()

	return CamelToSnake(structType.Name()), nil
}

func (*ReflectionHelper) ExtractTableFk(entity interface{}) (string, error) {
	v := reflect.ValueOf(entity)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, exception.
			New("Entity provided is not a struct pointer").
			Trace("reflect.ValueOf", "reflection.go")
	}

	structValue := v.Elem()
	structType := structValue.Type()

	scanTargets := make([]interface{}, structType.NumField())

	fieldMap := make(map[string]int)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		dbTag := field.Tag.Get("db")

		if dbTag == "" {
			continue
		}

		fieldMap[dbTag] = i

		scanTargets[i] = structValue.Field(i).Addr().Interface()
	}

	return scanTargets, nil
}
