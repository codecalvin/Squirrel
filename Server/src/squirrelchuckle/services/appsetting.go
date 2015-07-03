package services

import (
	"os"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"encoding/json"

	"squirrelchuckle/utility"
)

type AppSetting struct {
	appConfig map[string]string
	runConfig map[string]string

	AppDatabaseUrl string
	AppDatabasePort string
	AppDatabaseAddr string

	runConfigPath string
	alive bool
	dirty bool
}

var instance *AppSetting
func AppSettingInstance() *AppSetting {
	if instance == nil {
		instance = &AppSetting{}
	}
	return instance
}

func (this *AppSetting) RunConfig(key string) string{
	value, _ := this.runConfig[key]
	return value
}

func (this *AppSetting) UpdateRunConfig(key, value string) string {
	this.runConfig[key] = value
	this.dirty = true
	return value
}

func (this *AppSetting) Initialize() error {
	workPath, _ := os.Getwd()
	appPath, err := filepath.Abs(filepath.Dir(os.Args[0]))

	appConfigPath := filepath.Join(appPath, "conf", "appconf.json")
	this.runConfigPath = filepath.Join(appPath, "conf", "run.json")

	if workPath != appPath {
		if utility.FileExists(appConfigPath) {
			os.Chdir(appPath)
		} else {
			appConfigPath = filepath.Join(workPath, "conf", "appconf.json")
			this.runConfigPath = filepath.Join(workPath, "conf", "run.json")
		}
	}

	var ok bool
	settings, err := ioutil.ReadFile(appConfigPath)
	if err != nil {
		goto stale
	}

	this.appConfig = make(map[string]string)
	err = json.Unmarshal(settings, &this.appConfig)
	if err != nil {
		goto stale
	}

	this.AppDatabaseUrl, ok = this.appConfig["db_url"]
	if !ok {
		this.AppDatabaseUrl = "127.0.0.1"
	}

	this.AppDatabasePort, ok = this.appConfig["db_port"]
	if !ok {
		this.AppDatabasePort = "27017"
	}

	this.AppDatabaseAddr = fmt.Sprintf("%v:%v", this.AppDatabaseUrl, this.AppDatabasePort)

	this.runConfig = make(map[string]string)
	if utility.FileExists(this.runConfigPath) {
		settings, err = ioutil.ReadFile(this.runConfigPath)
		if err != nil {
			goto stale
		}
		err = json.Unmarshal(settings, &this.runConfig)
	}

	if err != nil {
		goto stale
	}

	this.alive = true

//success:
	return nil

stale:
	return err
}

func (this *AppSetting) UnInitialize() {
	if this.alive {
		this.Serialize()
		this.alive = false
	}
}

func (this *AppSetting) Alive() bool {
	return this.alive
}

func (this *AppSetting) Serialize() error {
	var err error
	if !this.alive || !this.dirty {
		return nil
	}

	var file *os.File
	content, err := json.MarshalIndent(this.runConfig, "", "    ")
	if err != nil {
		goto stale
	}

	file, err = os.Open(this.runConfigPath)
	if err != nil {
		goto stale
	}

	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		goto stale
	}

	this.dirty = false

//success:
	return nil

stale:
	return err
}