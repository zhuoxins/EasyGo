package maptranser

import "encoding/json"

func ToJson(maps interface{}) (string, error) {
	jsons, err := json.Marshal(maps)
	if err != nil {
		return "", err
	}
	return string(jsons), err
}
