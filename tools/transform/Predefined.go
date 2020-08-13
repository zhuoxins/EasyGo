package transform

func IsInt(vars interface{}) bool {
	if _, ok := vars.(int); ok {
		return true
	}
	return false
}

func IsString(vars interface{}) bool {
	if _, ok := vars.(string); ok {
		return true
	}
	return false
}

func IsInt64(vars interface{}) bool {
	if _, ok := vars.(int64); ok {
		return true
	}
	return false
}

func IsFloat(vars interface{}) bool {
	if _, ok := vars.(float64); ok {
		return true
	}
	if _, ok := vars.(float32); ok {
		return true
	}
	return false
}
