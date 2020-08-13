package builder

import (
	"EasyGo/tools/helper"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	ENTERFILE = "main"
	ROUTEFILE = "routes"
)

//创建项目
func NewProject() {
	var projectName string
	flag.StringVar(&projectName, "name", "", "new project name")
	flag.Parse()
	if projectName == "" {
		log.Println("project name must required")
		return
	}
	builder := &Builder{projectName: projectName, sep: sep}
	builder.Build()
}

type Builder struct {
	projectName string
	projectDir  string
	sep         string
	wg          sync.WaitGroup
}

func (this *Builder) Build() {
	//创建项目目录
	this.projectDir = helper.BasePath() + this.sep + this.projectName
	err := os.Mkdir(this.projectDir, 755)
	if err != nil {
		this.failed(err)
	}
	this.wg.Add(4)
	go this.createSysTemFile()
	go this.createRoutes()
	go this.createConf()
	go this.createApp()
	this.wg.Wait()
	log.Println("Project created successfully")
}

//创建系统文件
func (this *Builder) createSysTemFile() {
	//创建main入口文件
	mainPath := this.projectDir + this.sep + ENTERFILE + ".go"
	file, err := os.Create(mainPath)
	defer file.Close()
	if err != nil {
		this.failed(err)
		return
	}
	mainCon, _ := ioutil.ReadFile(itemTem(ENTERFILE))
	mainContent := strings.Replace(string(mainCon), PROJECTNAME, this.projectName, -1)
	mainContent = strings.Replace(mainContent, ROUTENAME, ROUTEFILE, -1)
	err = ioutil.WriteFile(mainPath, []byte(mainContent), 755)
	if err != nil {
		this.failed(err)
	}
	defer this.wg.Done()
	log.Println("system created successfully")
}

//创建路由文件
func (this *Builder) createRoutes() {
	routeDir := this.projectDir + this.sep + ROUTEFILE
	os.Mkdir(routeDir, 755)
	routeFile := routeDir + this.sep + ROUTEFILE + ".go"
	route, err := os.Create(routeFile)
	defer route.Close()
	if err != nil {
		this.failed(err)
		return
	}
	routeCon, _ := ioutil.ReadFile(itemTem(ROUTEFILE))
	routeContent := strings.Replace(string(routeCon), PACKAGENAME, ROUTEFILE, -1)
	err = ioutil.WriteFile(routeFile, []byte(routeContent), 755)
	if err != nil {
		this.failed(err)
	}
	defer this.wg.Done()
	log.Println("route created successfully")
}

//创建配置文件目录
func (this *Builder) createConf() {
	confDir := this.projectDir + this.sep + "config"
	_ = os.Mkdir(confDir, 755)
	//创建app.conf
	app, _ := os.Create(confDir + this.sep + "app.conf")
	defer app.Close()
	_ = ioutil.WriteFile(confDir+this.sep+"app.conf", []byte(baseConf("app")), 755)
	//创建cache.conf
	cache, _ := os.Create(confDir + this.sep + "cache.conf")
	defer cache.Close()
	_ = ioutil.WriteFile(confDir+this.sep+"cache.conf", []byte(baseConf("cache")), 755)
	//创建database.conf
	database, _ := os.Create(confDir + this.sep + "database.conf")
	defer database.Close()
	_ = ioutil.WriteFile(confDir+this.sep+"database.conf", []byte(baseConf("database")), 755)
	defer this.wg.Done()
	log.Println("config created successfully")
}

func (this *Builder) createApp() {
	appDir := this.projectDir + this.sep + "app"
	_ = os.Mkdir(appDir, 755)
	controllerDir := appDir + this.sep + "controllers"
	_ = os.Mkdir(controllerDir, 755)
	NewController("Index", controllerDir).build()
	middleDir := appDir + this.sep + "middleware"
	_ = os.Mkdir(middleDir, 755)
	modelDir := appDir + this.sep + "models"
	_ = os.Mkdir(modelDir, 755)
	NewModel("User", modelDir).build()
	defer this.wg.Done()
}

func (this *Builder) failed(err error) {
	if this.projectName != "" {
		go this.reset()
	}
	log.Println("create failed:", err.Error())
}

func (this *Builder) reset() {
	_ = os.RemoveAll(this.projectDir)
}
