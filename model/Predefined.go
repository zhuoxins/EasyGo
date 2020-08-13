package model

import (
	"reflect"
	"strings"
)

func New(model interface{}) {
	instance := reflect.ValueOf(model)
	tabName := instance.Elem().FieldByName("TableName").String()
	if tabName == "" {
		tabName = reflect.TypeOf(model).String()
		if strings.Contains(tabName, ".") {
			parseTabName := strings.Split(tabName, ".")
			tabName = parseTabName[len(parseTabName)-1]
		}
		tabName = strings.ToLower(tabName)
		if strings.Contains(tabName, "model") {
			tabName = strings.TrimRight(tabName, "model")
		}
		instance.Elem().FieldByName("TableName").Set(reflect.ValueOf(tabName))
	}
	instance.MethodByName("Init").Call(make([]reflect.Value, 0))
}
