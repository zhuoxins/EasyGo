package builder

import (
	"EasyGo/tools/helper"
	"os"
)

const (
	PACKAGENAME = "PACKAGENAME"
	CONTROLNAME = "CONTROLLERNAME"
	MODELNAME   = "MODELNAME"
	MIDDLENAME  = "MIDDLEWARENAME"
	PROJECTNAME = "PROJECTNAME"
	ROUTENAME   = "ROUTENAME"
	APP_DIR     = "app"
)

var sep = helper.Sep()

func controllerDir() (string, error) {
	controllerPath := helper.Pwd() + sep + APP_DIR + sep + "controllers"
	if !helper.PathExists(controllerPath) {
		//不存在则初始化目录
		err := os.MkdirAll(controllerPath, 755)
		return controllerPath, err
	}
	return controllerPath, nil
}

func modelDir() (string, error) {
	modelPath := helper.Pwd() + sep + APP_DIR + sep + "models"
	if !helper.PathExists(modelPath) {
		//不存在则初始化目录
		err := os.MkdirAll(modelPath, 755)
		return modelPath, err
	}
	return modelPath, nil
}

func middleDir() (string, error) {
	modelPath := helper.Pwd() + sep + APP_DIR + sep + "middleware"
	if !helper.PathExists(modelPath) {
		//不存在则初始化目录
		err := os.MkdirAll(modelPath, 755)
		return modelPath, err
	}
	return modelPath, nil
}

func getTemFilePath() string {
	path := helper.Pwd() + sep + "tools" + sep + "builder" + sep + "fileTem" + sep
	return path
}

func itemTem(file ...string) string {
	temDir := helper.Pwd() + sep + "tools" + sep + "builder" + sep + "projectTem"
	if len(file) == 0 {
		return temDir + sep
	}
	return temDir + sep + file[0] + ".base"
}

func controllerTem() string {
	return `package PACKAGENAME

import "EasyGo/route"

type CONTROLLERNAMEController struct {
    route.Controller
}

func (this *CONTROLLERNAMEController) Index() {
	//
}`
}

func modelTem() string {
	return `package PACKAGENAME

import "EasyGo/model"

type MODELNAMEModel struct {
	model.Model
	tabName string
}

//instancing model
func NewMODELNAME() *MODELNAMEModel {
	return model.New(&MODELNAMEModel{tabName: "MODELNAME"}).(*MODELNAMEModel)
}`
}

func middleTem() string {
	return `package PACKAGENAME

import (
	"EasyGo/Context"
)

func MIDDLEWARENAME(ctx *Context.Context) {
	//
}`
}

func baseConf(types string) string {
	switch types {
	case "app":
		return `[basis]
port = 8421
timezone = PRC`
	case "cache":
		return `driver =
prefix =
[redis]
network = tcp
host = 127.0.0.1
port = 6379
auth = 
`
	case "database":
		return `
connection = mysql
[mysql]
host = 127.0.0.1
database = blog
username = root
password = 
charset = utf8
`
	default:
		return ""
	}
}
