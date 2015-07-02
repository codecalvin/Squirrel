package settings

import (
	"os"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"encoding/json"

	"squirrelchuckle/utility"
)

var AppConfig map[string]string
var RunConfig map[string]string

var AppDatabaseUrl string
var AppDatabasePort string
var AppDatabaseAddr string

var runConfigPath string
var alive bool

func init() {
	workPath, _ := os.Getwd()
	workPath, _ = filepath.Abs(workPath)
	// initialize default configurations
	appPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	appConfigPath := filepath.Join(appPath, "conf", "appconf.json")
	runConfigPath = filepath.Join(appPath, "conf", "run.json")
	
	if workPath != appPath {
		if utility.FileExists(appConfigPath) {
			os.Chdir(appPath)
		} else {
			appConfigPath = filepath.Join(workPath, "conf", "appconf.json")
			runConfigPath = filepath.Join(workPath, "conf", "run.json")
		}
	}

	settings, err := ioutil.ReadFile(appConfigPath)
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
	
	RunConfig = make(map[string]string)
	if utility.FileExists(runConfigPath) {
		settings, err = ioutil.ReadFile(runConfigPath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(settings, &RunConfig)
	}
	
	if err != nil {
		panic(err)
	}
	
	alive = true
}

func Serialize() {
	if alive {
		content, err := json.MarshalIndent(RunConfig, "", "    ")
		if err != nil {
			fmt.Errorf("[FATAL ERROR] failed to output run config, need rescue")
			return
		}

		file, err := os.Open(runConfigPath)
		if err != nil {
			fmt.Errorf("[FATAL ERROR] failed to output run config, need rescue. Cannot open file %v", runConfigPath)
			return
		}
		defer file.Close()
		file.Write(content)
	}
}