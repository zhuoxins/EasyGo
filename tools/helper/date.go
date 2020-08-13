package helper

import (
	"strings"
	"time"
)

var format map[string]string

func init() {
	format = make(map[string]string)
	format = map[string]string{"Y": "2006", "M": "01", "D": "02", "H": "15", "I": "04", "S": "05"}
}

//格式化时间戳
func Date(args ...interface{}) string {
	var dateFormat string
	var timestamp int64
	lenParse := len(args)
	if lenParse == 0 {
		dateFormat = "2006-01-02 15:04:05"
	} else {
		formatStr := strings.Split(args[0].(string), "")
		for _, v := range formatStr {
			if s, ok := format[strings.ToUpper(v)]; ok {
				dateFormat += s
			} else {
				dateFormat += v
			}
		}
	}
	if lenParse == 2 {
		timestamp = args[1].(int64)
	} else {
		loc, _ := time.LoadLocation("PRC")
		timestamp = time.Now().In(loc).Unix()
	}
	return time.Unix(timestamp, 0).Format(dateFormat)
}

func Strtotime(parseDate string, args ...interface{}) int64 {
	var dateFormat string
	var loc *time.Location
	lenParse := len(args)
	if lenParse == 0 {
		dateFormat = "2006-01-02 15:04:05"
	} else {
		formatStr := strings.Split(args[0].(string), "")
		for _, v := range formatStr {
			if s, ok := format[strings.ToUpper(v)]; ok {
				dateFormat += s
			} else {
				dateFormat += v
			}
		}
	}
	if lenParse == 2 {
		loc = args[1].(*time.Location)
	} else {
		loc, _ = time.LoadLocation("PRC")
	}
	t, err := time.ParseInLocation(dateFormat, parseDate, loc)
	if err != nil {
		return 0
	}
	return t.Unix()
}
