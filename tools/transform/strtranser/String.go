package strtranser

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strconv"
)

func ToInt(item string) int {
	i, err := strconv.Atoi(item)
	if err != nil {
		return 0
	}
	return i
}

func ToInt64(item string) int64 {
	i, err := strconv.ParseInt(item, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func JsonMap(str string, data interface{}, isFilter ...bool) error {
	if len(isFilter) > 0 {
		check := isFilter[0]
		if check {
			return json.Unmarshal(bytes.TrimPrefix([]byte(str), []byte("\xef\xbb\xbf")), data)
		}
	}
	return json.Unmarshal([]byte(str), data)
}
