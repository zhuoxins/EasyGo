package config

import (
	"EasyGo/kernel"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

var handler *Config

func GetConf(item string) interface{} {
	return handler.Get(item)
}

func Get(item string) Result {
	res := handler.Get(item)
	result := &result{value: res}
	return result
}

type Config struct {
	dirName    string
	configTree map[string]interface{}
}

func NewConfig(name string) error {
	//获取配置文件目录
	name = path(name)
	config := &Config{
		dirName: name,
	}
	err := config.load()
	if err != nil {
		return err
	}
	handler = config
	return nil
}

func (this *Config) Get(item string) interface{} {
	parseField := strings.Split(item, ".")
	if len(this.configTree) == 0 || this.configTree == nil {
		return nil
	}
	switch length := len(parseField); length {
	case 1:
		result, ok := this.configTree[parseField[0]]
		if ok {
			return result
		} else {
			return nil
		}
	case 2:
		if con, ok := this.configTree[parseField[0]]; ok {
			if conf, ok := con.(map[string]interface{}); ok {
				if result, ok := conf[parseField[1]]; ok {
					return result
				}
			}
		}
		return nil
	case 3:
		//获取某一模块下一字段值
		if con, ok := this.configTree[parseField[0]]; ok {
			if conChild, ok := con.(map[string]interface{})[parseField[1]]; ok {
				if res, ok := conChild.(map[string]string)[parseField[2]]; ok {
					return res
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func (this *Config) load() error {
	confContent, err := this.readConfig()
	if err != nil {
		return err
	}
	this.configTree = this.format(confContent)
	return nil
}

func (this *Config) readConfig() (map[string][]string, error) {
	confContents := make(map[string][]string)
	//遍历配置目录下所有配置文件
	entries, err := ioutil.ReadDir(this.dirName)
	if err != nil {
		return nil, err
	}
	for _, file := range entries {
		if !file.IsDir() && strings.Contains(file.Name(), ".conf") {
			confName := strings.TrimRight(file.Name(), ".conf")
			content, _ := ioutil.ReadFile(this.dirName + file.Name())
			con := strings.Split(string(content), "\n")
			confContents[confName] = con
		}
	}
	return confContents, nil
}

func (this *Config) format(contents map[string][]string) map[string]interface{} {
	result := make(map[string]interface{})
	var sectionName string
	for confName, con := range contents {
		result[confName] = make(map[string]interface{})
		sectionName = ""
		for _, lineStr := range con {
			lineStr = strings.TrimSpace(lineStr)
			//忽略空行和注释
			if lineStr == "" || lineStr[0] == ';' {
				continue
			}
			//行首和尾巴分别是方括号的说明是段标记的起止符
			if lineStr[0] == '[' && lineStr[len(lineStr)-1] == ']' {
				sectionName = lineStr[1 : len(lineStr)-1]
				result[confName].(map[string]interface{})[sectionName] = make(map[string]string)
			} else {
				//切开等号分割的键值对
				pair := strings.Split(lineStr, "=")
				//保证切开只有一个等号情况
				if len(pair) == 2 {
					key := strings.Trim(strings.TrimSpace(pair[0]), `"`)
					val := strings.Trim(strings.TrimSpace(pair[1]), `"`)
					if sectionName != "" && reflect.TypeOf(result[confName].(map[string]interface{})[sectionName]).String() == "map[string]string" {
						result[confName].(map[string]interface{})[sectionName].(map[string]string)[key] = val
					} else {
						result[confName].(map[string]interface{})[key] = val
					}
				}
			}
		}
	}
	return result
}

func init() {
	kernel.Register("Config").Do(func() {
		err := NewConfig("config")
		if err != nil {
			fmt.Println(err)
		}
	})
}
