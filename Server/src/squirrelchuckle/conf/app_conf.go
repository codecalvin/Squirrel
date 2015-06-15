package conf

import (
	"os"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"encoding/json"

	"squirrelchuckle/utility"
)

var AppConfig map[string]string
var AppDatabaseUrl string
var AppDatabasePort string
var AppDatabaseAddr string

func init() {
	workPath, _ := os.Getwd()
	workPath, _ = filepath.Abs(workPath)
	// initialize default configurations
	appPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	appInConfigPath := filepath.Join(appPath, "conf", "inapp.conf")
	
	if workPath != appPath {
		if utility.FileExists(appInConfigPath) {
			os.Chdir(appPath)
		} else {
			appInConfigPath = filepath.Join(workPath, "conf", "inapp.conf")
		}
	}

	settings, err := ioutil.ReadFile(appInConfigPath)
	if err != nil {
		panic(err)
	}
	
	AppConfig = make(map[string]string)
	err = json.Unmarshal(settings, &AppConfig)
	if err != nil {
		panic(err)
	}

	var ok bool
	AppDatabaseUrl, ok = AppConfig["db_url"]
	if !ok {
		AppDatabaseUrl = "127.0.0.1"
	}

	AppDatabasePort, ok = AppConfig["db_port"]
	if !ok {
		AppDatabasePort = "27017"
	}

	AppDatabaseAddr = fmt.Sprintf("%v:%v", AppDatabaseUrl, AppDatabasePort)
}