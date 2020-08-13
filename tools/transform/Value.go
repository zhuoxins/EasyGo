package transform

import (
	"EasyGo/tools/transform/numtranser"
	"EasyGo/tools/transform/strtranser"
)

type Value interface {
	ToInt() int
	ToInt64() int64
	ToString() string
}

func Transform(value interface{}) Value {
	return &result{handle: value}
}

type result struct {
	handle interface{}
}

func (r *result) ToInt() int {
	if IsInt(r.handle) {
		return r.handle.(int)
	} else if IsInt64(r.handle) {
		return numtranser.ToInt(r.handle)
	} else if IsString(r.handle) {
		return strtranser.ToInt(r.handle.(string))
	} else if IsFloat(r.handle) {
		return numtranser.ToInt(r.handle)
	}
	return 0
}

func (r *result) ToInt64() int64 {
	if IsInt(r.handle) {
		return numtranser.ToInt64(r.handle)
	} else if IsInt64(r.handle) {
		return r.handle.(int64)
	} else if IsString(r.handle) {
		return strtranser.ToInt64(r.handle.(string))
	} else if IsFloat(r.handle) {
		return numtranser.ToInt64(r.handle)
	}
	return 0
}

func (r *result) ToString() string {
	if IsInt(r.handle) {
		return numtranser.ToString(r.handle)
	} else if IsInt64(r.handle) {
		return numtranser.ToString(r.handle)
	} else if IsString(r.handle) {
		return r.handle.(string)
	} else if IsFloat(r.handle) {
		return numtranser.ToString(r.handle)
	}
	return ""
}
