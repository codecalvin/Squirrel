package core

import (
	"os"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"encoding/json"

	"squirrelchuckle/utility"
	"strconv"
)

type AppSetting struct {
	appConfig 		map[string]string
	runConfig 		map[string]string

	AppDatabaseUrl 	string
	AppDatabasePort string
	AppDatabaseAddr string

	// exchange service
	ExchangeAuth  	bool
	ExchangeHost  	string
	ExchangeUrl   	string
	ExchangePort  	uint

	runConfigPath 	string
	alive 			bool
	dirty		 	bool
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

func (this *AppSetting) Depends() []string {
	return nil
}

func (this *AppSetting) Alive() bool {
	return this.alive
}

func (this *AppSetting) Name() string {
	return "AppSetting"
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

	ok := true
	settings, err := ioutil.ReadFile(appConfigPath)
	if err != nil {
		goto stale
	}

	this.appConfig = make(map[string]string)
	if err = json.Unmarshal(settings, &this.appConfig); err != nil {
		goto stale
	}

	if this.AppDatabaseUrl, ok = this.appConfig["db_url"]; !ok {
		this.AppDatabaseUrl = "127.0.0.1"
	}
	if this.AppDatabasePort, ok = this.appConfig["db_port"]; !ok {
		this.AppDatabasePort = "27017"
	}
	this.AppDatabaseAddr = fmt.Sprintf("%v:%v", this.AppDatabaseUrl, this.AppDatabasePort)

	// exchange service setting
	if value, ok := this.appConfig["exchange_enable"]; ok {
		if this.ExchangeAuth, _ = strconv.ParseBool(value); this.ExchangeAuth {
			if this.ExchangeHost, ok = this.appConfig["exchange_server"]; !ok {
				this.ExchangeAuth = false
			}
			this.ExchangePort = 597
			if value, ok = this.appConfig["exchange_port"]; ok {
				if value, err := strconv.ParseUint(value, 10, 0); err != nil {
					SquirrelApp.Error("[Appsetting] error 'exchange_port'")
				} else {
					this.ExchangePort = uint(value)
				}
			}
			this.ExchangeUrl = fmt.Sprintf("%v:%v", this.ExchangeHost, this.ExchangePort)
		}
	}

	if this.runConfig = make(map[string]string); utility.FileExists(this.runConfigPath) {
		if settings, err = ioutil.ReadFile(this.runConfigPath); err == nil {
			err = json.Unmarshal(settings, &this.runConfig)
		}
	}

	if err != nil {
		goto stale
	}

	this.alive = true

//success:
	return nil

stale:
	SquirrelApp.Critical("AppSetting Initialize failed. Error: %v", err)
	return err
}

func (this *AppSetting) UnInitialize() {
	if this.alive {
		this.Serialize()
		this.alive = false
	}
}

func (this *AppSetting) Serialize() error {
	var err error
	if !this.alive || !this.dirty {
		return nil
	}

	var file *os.File
	var content []byte
	if content, err = json.MarshalIndent(this.runConfig, "", "    "); err != nil {
		goto stale
	}

	if file, err = os.Open(this.runConfigPath); err != nil {
		goto stale
	}

	defer file.Close()
	if _, err = file.Write(content); err != nil {
		goto stale
	}

	this.dirty = false

//success:
	return nil

stale:
	SquirrelApp.Critical("AppSetting serialize failed, need rescue. Error: %v", err)
	return err
}