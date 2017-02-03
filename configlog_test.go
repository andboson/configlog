package configlog

import (
	"log"
	"os"
	"regexp"
	"testing"
)

var testConfig = `
logfile: "logfile.log"
debug: true
redis:
  port: 6379`

func TestNotFoundConfig(t *testing.T) {
	if AppConfig != nil {
		t.Fatalf("AppConf is not nil!: %v", AppConfig)
	}
}

func TestConfigFound(t *testing.T) {
	os.Mkdir("config", 0777)
	f, err := os.OpenFile("config/app.yml", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}

	f.WriteString(testConfig)
	f.Close()

	//load conf
	load()
	logfilename, _ := AppConfig.String("logfile")
	if len(logfilename) == 0 {
		t.Errorf("Cannot read config %s", AppConfig)
	}

	//read value
	redisport, _ := AppConfig.String("redis.port")
	if redisport != "6379" {
		t.Errorf("Error reading config value! %s", redisport)
	}

	//check logs
	var content = make([]byte, 1000)
	log.Println("testlogentry")
	logfile, err := os.OpenFile("logfile.log", os.O_RDONLY, 0666)
	logfile.Read(content)
	contentString := string(content)
	match, _ := regexp.MatchString("testlogentry", contentString)
	if !match {
		t.Fatalf("Unable to find test log entry \n = %s", content)
	}

	//clearup
	os.Remove("logfile.log")
	os.Remove("config/app.yml")
	os.Remove("config")
}
