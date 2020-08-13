package view

import (
	"EasyGo/context"
	"EasyGo/tools/helper"
	"html/template"
	"io/ioutil"
	"path"
)

//模板加载核心文件
var templates = make(map[string]*template.Template)

//获取视图目录路径
func getViewPath() string {
	//项目路径
	dirPath := helper.Pwd()
	//目录分隔符
	sep := helper.Sep()
	//模板文件目录
	viewPath := dirPath + sep + "app" + sep + "views"
	return viewPath
}

//初始化缓存静态模板
func InitViewResources() {
	viewPath := getViewPath()
	sep := helper.Sep()
	fileArr, err := ioutil.ReadDir(viewPath)
	if err != nil {
		panic(err)
		return
	}
	var templateName, templatePath string
	for _, fileInfo := range fileArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = viewPath + sep + templateName
		//log.Println("Loading template ...")
		t := template.Must(template.ParseFiles(templatePath))
		templates[templateName] = t
	}
}

//加载视图
func View(ctx *context.Context, fileName string, data map[string]interface{}) (err error) {
	//拼接模板文件名
	fileName = fileName + ".html"
	//加载模板文件
	err = templates[fileName].Execute(ctx.Response.Response(), data)
	if err != nil {
		ctx.Response.Error(err.Error())
	}
	return
}
