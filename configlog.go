package configlog

import (
	"github.com/olebedev/config"
	"path/filepath"
	 "os"
	"io/ioutil"
	"regexp"
	"strings"
	"github.com/kardianos/osext"
	log "github.com/Sirupsen/logrus"
)

const (
	CONFIG_DIR = "config"
	PRODUCTION_FOLDER = "production"
	CONFIG_FILE = "app.yml"
)

var AppConfig *config.Config
var CurrDirectory string;

func init(){
	load()
}

func load(){
	var err error
	var yml []byte
	configFile := detectProdConfig(false)
	yml, err = ioutil.ReadFile(configFile)
	if(err != nil ){
		log.Printf("Unable to find config in path: %s,  %s", configFile, err)
		return
	}
	AppConfig, err = config.ParseYaml(string(yml))
	log.SetFormatter(&log.JSONFormatter{})
	EnableLogfile()
}

func detectProdConfig(useosxt bool) string{
	var levelUp string
	var curDir string
	sep := string(filepath.Separator)

	if(useosxt){
		curDir, _ = os.Getwd()
	}else {
		curDir, _ = osext.ExecutableFolder()
	}

	//detect from test or console
	match, _ := regexp.MatchString("_test",curDir)
	matchArgs, _ := regexp.MatchString("arguments",curDir)
	matchTestsDir, _ := regexp.MatchString("tests",curDir)
	if(match || matchArgs || matchTestsDir){
		if(matchTestsDir){
			levelUp = ".."
		}
		curDir, _ = filepath.Abs(curDir + sep+ levelUp + sep)
	}

	CurrDirectory = curDir;
	configDir, _ := filepath.Abs(curDir + sep + CONFIG_DIR + sep)
	appConfig := configDir + sep + CONFIG_FILE
	appProdConfig := configDir + sep + PRODUCTION_FOLDER + sep + CONFIG_FILE
	if(fileExists(appProdConfig)){
		appConfig = appProdConfig
	} else if(!useosxt){
		appConfig = detectProdConfig(true)
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