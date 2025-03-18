package utils

import (
	"reflect"
	"sync"
)

var modelFieldCache = sync.Map{}

func RegisterModels(models ...interface{}) {
	for _, model := range models {
		registerModel(model)
	}
}

func registerModel(model interface{}) {
	t := reflect.TypeOf(model)

	fieldMap := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		column := field.Name
		if tag, ok := field.Tag.Lookup("db"); ok {
			column = tag
		}
		fieldMap[field.Name] = column
	}

	modelFieldCache.Store(t.Name(), fieldMap)
}

func GetField(model interface{}, fieldName string) string {
	t := reflect.TypeOf(model)
	if fields, ok := modelFieldCache.Load(t.Name()); ok {
		if fieldMap, ok := fields.(map[string]string); ok {
			if column, exists := fieldMap[fieldName]; exists {
				return column
			}
		}
	}
	panic("Field not found: " + fieldName)
}
