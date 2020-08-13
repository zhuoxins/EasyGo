package config

import (
	"EasyGo/tools/transform"
	"EasyGo/tools/transform/numtranser"
	"EasyGo/tools/transform/strtranser"
)

//结果类型转换
type result struct {
	value interface{}
}

func (this *result) String(defaultVal ...string) string {
	var res string
	if transform.IsString(this.value) {
		res = this.value.(string)
	} else if transform.IsInt(this.value) {
		res = numtranser.ToString(this.value)
	} else if transform.IsFloat(this.value) {
		res = numtranser.ToString(this.value)
	}
	if res == "" && len(defaultVal) > 0 {
		defVal := defaultVal[0]
		return defVal
	}
	return res
}

func (this *result) Int() int {
	var i int
	if transform.IsInt(this.value) {
		i = this.value.(int)
	} else if transform.IsString(this.value) {
		i = strtranser.ToInt(this.value.(string))
	} else if transform.IsInt64(this.value) {
		i = numtranser.ToInt(this.value)
	} else if transform.IsFloat(this.value) {
		i = numtranser.ToInt(this.value)
	}
	return i
}

func (this *result) Int64() int64 {
	var i int64
	if transform.IsInt(this.value) {
		i = numtranser.ToInt64(this.value)
	} else if transform.IsString(this.value) {
		i = strtranser.ToInt64(this.value.(string))
	} else if transform.IsInt64(this.value) {
		i = this.value.(int64)
	} else if transform.IsFloat(this.value) {
		i = numtranser.ToInt64(this.value)
	}
	return i
}

func (this *result) Map() map[string]Result {
	data := make(map[string]Result)
	if stringMap, ok := this.value.(map[string]string); ok {
		for key, val := range stringMap {
			data[key] = &result{value: val}
		}
	} else if maps, ok := this.value.(map[string]interface{}); ok {
		for key, val := range maps {
			data[key] = &result{value: val}
		}
	}
	return data
}

func (this *result) Item(field string, defaultVal ...string) string {
	var defVal string
	if len(defaultVal) > 0 {
		defVal = defaultVal[0]
	}
	if stringMap, ok := this.value.(map[string]string); ok {
		if val, ok := stringMap[field]; ok {
			return val
		}
	} else if maps, ok := this.value.(map[string]interface{}); ok {
		if val, ok := maps[field]; ok {
			if transform.IsString(val) {
				return val.(string)
			} else if transform.IsInt(val) {
				return numtranser.ToString(val)
			} else if transform.IsFloat(val) {
				return numtranser.ToString(this.value)
			}
		}
	}
	return defVal
}

func (this *result) GetField(field string, defaultVal ...interface{}) Result {
	r := &result{}
	var defVal interface{}
	if len(defaultVal) > 0 {
		defVal = defaultVal[0]
	}
	if stringMap, ok := this.value.(map[string]string); ok {
		if val, ok := stringMap[field]; ok {
			r.value = val
		} else {
			if defVal != nil {
				r.value = defVal
			}
		}
	} else if maps, ok := this.value.(map[string]interface{}); ok {
		if val, ok := maps[field]; ok {
			r.value = val
		} else {
			if defVal != nil {
				r.value = defVal
			}
		}
	}
	if r.value == nil {
		return nil
	}
	return r
}
