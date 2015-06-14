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

func init() {
	workPath, _ := os.Getwd()
	workPath, _ = filepath.Abs(workPath)
	// initialize default configurations
	appPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	appInConfigPath := filepath.Join(appPath, "conf", "inapp.conf")

	if !utility.FileExists(appInConfigPath) {

	}

	settings, err := ioutil.ReadFile(appInConfigPath)
	if err != nil {
		panic(fmt.Sprintf("Cannot read %v", appInConfigPath))
	}

	json.Unmarshal(settings, AppConfig)
}