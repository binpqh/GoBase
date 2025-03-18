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

func GetFields[T any]() map[string]string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if fields, ok := modelFieldCache.Load(t.Name()); ok {
		if fieldMap, ok := fields.(map[string]string); ok {
			return fieldMap
		}
	}
	panic("Fields not found for model: " + t.Name())
}
