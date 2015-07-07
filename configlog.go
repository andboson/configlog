package configlog

import (
	"github.com/olebedev/config"
	"path/filepath"
	"os"
	"log"
	"io/ioutil"
	"runtime"
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
	curDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	//detect from test or console
	match, _ := regexp.MatchString("_test",curDir)
	matchArgs, _ := regexp.MatchString("arguments",curDir)
	if(match || matchArgs){
		matchTestsDir, _ := regexp.MatchString("tests",curDir)
		if(matchTestsDir){
			levelUp = ".."
		}
		_, file, _, _ := runtime.Caller(1)
		callerPath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, levelUp + string(filepath.Separator))))
		curDir = callerPath
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