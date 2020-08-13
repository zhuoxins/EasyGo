package builder

import (
	"io/ioutil"
	"os"
	"strings"
)

type ControllerBuilder struct {
	name string
	dir  string
}

func NewController(name string, dir ...string) *ControllerBuilder {
	builder := &ControllerBuilder{name: name}
	if len(dir) > 0 {
		builder.dir = dir[0]
	}
	return builder
}

func (this *ControllerBuilder) build() error {
	//初始化构建目录
	if this.dir == "" {
		dir, err := controllerDir()
		if err != nil {
			return err
		}
		this.dir = dir
	}
	packName := "controllers"
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
	//创建控制器文件
	controllerFile := this.dir + sep + this.name + ".go"
	file, err := os.Create(controllerFile)
	if err != nil {
		return err
	}
	defer file.Close()
	content := controllerTem()
	content = strings.Replace(content, "PACKAGENAME", packName, -1)
	content = strings.Replace(content, "CONTROLLERNAME", this.name, -1)
	err = ioutil.WriteFile(controllerFile, []byte(content), 755)
	if err != nil {
		_ = os.Remove(controllerFile)
		return err
	}
	return nil
}
