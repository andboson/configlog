package configlog

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/howeyc/fsnotify"
	"github.com/kardianos/osext"
	"github.com/olebedev/config"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
)

const (
	CONFIG_DIR        = "config"
	PRODUCTION_FOLDER = "production"
	CONFIG_FILE       = "app.yml"
)

var AppConfig *config.Config
var CurrDirectory string
var m sync.RWMutex
var Out *os.File

func init() {
	ReloadConfigLog()
	watchLog()
}

func watchLog() {
	if AppConfig == nil {
		return
	}
	sigs := make(chan os.Signal, 1)

	logfileName, _ := AppConfig.String("logfile")
	watcher, err := fsnotify.NewWatcher()
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGTERM)

	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsRename() || ev.IsDelete() {
					ReloadConfigLog()
				}
			case err := <-watcher.Error:
				log.Println("[configlog]  log watcher error:", err)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-sigs:
				ReloadConfigLog()
			}
		}

	}()
	err = watcher.Watch(logfileName)
}

func ReloadConfigLog() {
	m.Lock()
	defer m.Unlock()
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{})
	load()
	log.Printf("[configlog] reaload config")

	return
}

func load() {
	var err error
	var yml []byte
	configFile := detectProdConfig(false)
	yml, err = ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("[configlog]  Unable to find config in path: %s,  %s", configFile, err)
		return
	}
	AppConfig, err = config.ParseYaml(string(yml))
	logfileName, _ := AppConfig.String("logfile")
	Out = EnableLogfile(logfileName)
}

func detectProdConfig(useosxt bool) string {
	var levelUp string
	var curDir string
	sep := string(filepath.Separator)

	if useosxt {
		curDir, _ = os.Getwd()
	} else {
		curDir, _ = osext.ExecutableFolder()
	}

	//detect from test or console
	match, _ := regexp.MatchString("_test", curDir)
	matchArgs, _ := regexp.MatchString("arguments", curDir)
	matchTestsDir, _ := regexp.MatchString("tests", curDir)
	if match || matchArgs || matchTestsDir {
		if matchTestsDir {
			levelUp = ".."
		}
		curDir, _ = filepath.Abs(curDir + sep + levelUp + sep)
	}

	CurrDirectory = curDir
	configDir, _ := filepath.Abs(curDir + sep + CONFIG_DIR + sep)
	appConfig := configDir + sep + CONFIG_FILE
	appProdConfig := configDir + sep + PRODUCTION_FOLDER + sep + CONFIG_FILE
	if fileExists(appProdConfig) {
		appConfig = appProdConfig
	} else if !useosxt {
		appConfig = detectProdConfig(true)
	}

	return appConfig
}

func EnableLogfile(logfileName string) *os.File {

	if logfileName == "" {
		log.Printf("[configlog]  logfile is STDOUT")
		return nil
	}

	log.Printf("[configlog]  logfile is %s", logfileName)
	logFile := logfileName
	logfileNameSlice := strings.Split(logfileName, string(filepath.Separator))

	//relative path
	if len(logfileNameSlice) > 1 && logfileNameSlice[0] != "" {
		logFile = CurrDirectory + string(filepath.Separator) + logfileName
	}

	//try to create log folder
	if len(logfileNameSlice) > 1 {
		logfileNameSlice = logfileNameSlice[:len(logfileNameSlice)-1]
		logPath := strings.Join(logfileNameSlice, string(filepath.Separator))
		os.Mkdir(logPath, 0777)
	}

	if Out != nil {
		log.Printf("[configlog] closing file")
		Out.Close()
	}
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[configlog] error opening file: %v", err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(f)
	Out = f

	return f
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
