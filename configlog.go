package configlog

import (
	"github.com/olebedev/config"
	"path/filepath"
	 "os"
	"log"
	"io/ioutil"
	"regexp"
	"strings"
)

var AppConfig *config.Config
var CurrDirectory string;

func init(){
	load()
}

func load(){
	configFile := detectProdConfig()
	yml, err := ioutil.ReadFile(configFile)
	AppConfig, err = config.ParseYaml(string(yml))
	if(err != nil){
		log.Printf("Unable to find config in path: %s,  %s", configFile, err)
		return
	}
	EnableLogfile()
}

func detectProdConfig() string{
	var levelUp string
	sep := string(filepath.Separator)
	curDir, _ := os.Getwd()

	//detect from test or console
	match, _ := regexp.MatchString("_test",curDir)
	matchArgs, _ := regexp.MatchString("arguments",curDir)
	matchTestsDir, _ := regexp.MatchString("tests",curDir)
	if(match || matchArgs || matchTestsDir){
		if(matchTestsDir){
			levelUp = ".."
		}
		curDir, _ = filepath.Abs(curDir + string(filepath.Separator) + levelUp + string(filepath.Separator))
	}

	CurrDirectory = curDir;
	configDir, _ := filepath.Abs(curDir + sep +"config" + sep)
	appConfig := configDir + sep + "app.yml"
	appProdConfig := configDir + sep + "production" + sep + "app.yml"
	if(fileExists(appProdConfig)){
		appConfig = appProdConfig
	}

	return appConfig
}

func EnableLogfile(){
	logfileName, _ := AppConfig.String("logfile")

	if(logfileName == ""){
		log.Printf("logfile is STDOUT")
		return;
	}

	log.Printf("logfile is %s", logfileName)
	logFile := logfileName
	logfileNameSlice := strings.Split(logfileName, string(filepath.Separator))

	//relative path
	if(len(logfileNameSlice) > 1 && logfileNameSlice[0] != ""){
		logFile = CurrDirectory + string(filepath.Separator)  +logfileName
	}

	//try to create log folder
	if(len(logfileNameSlice) > 1) {
		logfileNameSlice =  logfileNameSlice[:len(logfileNameSlice)-1]
		logPath := strings.Join(logfileNameSlice, string(filepath.Separator))
		os.Mkdir(logPath, 0777)
	}

	f, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}