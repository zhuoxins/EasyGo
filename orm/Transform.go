package orm

import (
	"EasyGo/tools/helper"
	"strconv"
)

type Value string
type ValueMap map[string]Value

func (v Value) ToInt() int {
	i, _ := strconv.Atoi(string(v))
	return i
}

func (v Value) ToInt64() int64 {
	i, _ := strconv.ParseInt(string(v), 10, 64)
	return i
}

func (v Value) MD5() string {
	return helper.MD5(string(v))
}

func Result(res string, err error) (Value, error) {
	return Value(res), err
}

func ResultMap(result interface{}, err error) (ValueMap, error) {
	var data ValueMap
	if res, ok := result.(map[string]string); ok {
		data = make(ValueMap)
		for key, val := range res {
			data[key] = Value(val)
		}
	}
	return data, err
}

func ResultMaps(result interface{}, err error) ([]ValueMap, error) {
	var data []ValueMap
	if res, ok := result.([]map[string]string); ok {
		for _, value := range res {
			dataMap := make(ValueMap)
			for key, val := range value {
				dataMap[key] = Value(val)
			}
			data = append(data, dataMap)
		}
	}
	return data, err
}
