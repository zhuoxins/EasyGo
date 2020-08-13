package helper

import (
	"EasyGo/kernel"
	"crypto/md5"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
)

func Sep() string {
	return string(os.PathSeparator)
}

func Pwd() string {
	path, _ := os.Getwd()
	return path
}

func EnvPath() string {
	return strings.TrimRight(Pwd(), Sep()+kernel.ProjectName)
}

//判断目录或文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func InMap(item string, maps map[string]interface{}) bool {
	if _, ok := maps[item]; ok {
		return true
	}
	return false
}

func TrimStr(str, leftStr, rightStr string) string {
	return strings.TrimLeft(strings.TrimRight(str, rightStr), leftStr)
}

func InSlice(key int, slices []string) bool {
	if (len(slices) - 1) < key {
		return false
	}
	return true
}

func InStringArray(field string, arr []string) bool {
	for _, item := range arr {
		if field == item {
			return true
		}
	}
	return false
}

func ToInt(item interface{}) (int, error) {
	if str, ok := item.(string); ok {
		return strconv.Atoi(str)
	}
	return 0, nil
}

func StrLength(str string) int {
	i := 0
	for range str {
		i++
	}
	return i
}

func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func UcFirst(str string) string {
	str = strings.ToLower(str)
	firstWord := str[0:1]
	str = strings.ToUpper(firstWord) + str[1:]
	return str
}

func BasePath() string {
	proDir := Pwd()
	return strings.TrimRight(proDir, Sep()+"EasyGo")
}
