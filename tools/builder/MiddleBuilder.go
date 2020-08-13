package builder

import (
	"io/ioutil"
	"os"
	"strings"
)

type MiddleBuilder struct {
	name string
	dir  string
}

func NewMiddle(name string, dir ...string) *MiddleBuilder {
	builder := &MiddleBuilder{name: name}
	if len(dir) > 0 {
		builder.dir = dir[0]
	}
	return builder
}

func (this *MiddleBuilder) build() error {
	//初始化构建目录
	if this.dir == "" {
		dir, err := middleDir()
		if err != nil {
			return err
		}
		this.dir = dir
	}
	packName := "middleware"
	//判断是否包含模块名
	if strings.Contains(this.name, ".") {
		parseName := strings.Split(this.name, ".")
		moduleName := parseName[0]
		this.name = parseName[1]
		packName = moduleName
		this.dir += sep + moduleName
		//创建模块目录
		err := os.Mkdir(this.dir, 755)
		if err != nil {
			return err
		}
	}
	//创建middleware文件
	middleFile := this.dir + sep + this.name + ".go"
	file, err := os.Create(middleFile)
	if err != nil {
		return err
	}
	defer file.Close()
	content := middleTem()
	content = strings.Replace(content, "PACKAGENAME", packName, -1)
	content = strings.Replace(content, "MIDDLEWARENAME", this.name, -1)
	err = ioutil.WriteFile(middleFile, []byte(content), 755)
	if err != nil {
		_ = os.Remove(middleFile)
		return err
	}
	return nil
}
