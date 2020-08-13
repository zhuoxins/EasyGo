package builder

import (
	"io/ioutil"
	"os"
	"strings"
)

type ModelBuilder struct {
	name string
	dir  string
}

func NewModel(name string, dir ...string) *ModelBuilder {
	builder := &ModelBuilder{name: name}
	if len(dir) > 0 {
		builder.dir = dir[0]
	}
	return builder
}

func (this *ModelBuilder) build() error {
	//初始化构建目录
	if this.dir == "" {
		dir, err := modelDir()
		if err != nil {
			return err
		}
		this.dir = dir
	}
	packName := "models"
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
	//创建model文件
	modelFile := this.dir + sep + this.name + ".go"
	file, err := os.Create(modelFile)
	if err != nil {
		return err
	}
	defer file.Close()
	content := modelTem()
	content = strings.Replace(content, "PACKAGENAME", packName, -1)
	content = strings.Replace(content, "MODELNAME", this.name, -1)
	err = ioutil.WriteFile(modelFile, []byte(content), 755)
	if err != nil {
		_ = os.Remove(modelFile)
		return err
	}
	return nil
}
