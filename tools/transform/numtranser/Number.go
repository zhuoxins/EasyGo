package numtranser

import "strconv"

func ToInt(item interface{}) (i int) {
	var err error
	if num, ok := item.(int64); ok {
		i, err = strconv.Atoi(strconv.FormatInt(num, 10))
	} else if fl, ok := item.(float64); ok {
		str := strconv.FormatFloat(fl, 'f', -1, 64)
		i, err = strconv.Atoi(str)
	}
	if err != nil {
		return 0
	}
	return i
}

func ToInt64(item interface{}) (i int64) {
	var err error
	if num, ok := item.(int); ok {
		str := strconv.Itoa(num)
		i, err = strconv.ParseInt(str, 10, 64)
	} else if fl, ok := item.(float64); ok {
		str := strconv.FormatFloat(fl, 'f', -1, 64)
		i, err = strconv.ParseInt(str, 10, 64)
	}
	if err != nil {
		return 0
	}
	return i
}

func ToString(item interface{}) (str string) {
	if num, ok := item.(int64); ok {
		str = strconv.FormatInt(num, 10)
	} else if fl, ok := item.(float64); ok {
		str = strconv.FormatFloat(fl, 'f', -1, 64)
	} else if num, ok := item.(int); ok {
		str = strconv.Itoa(num)
	}
	return
}
