package licenseload

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/choueric/clog"
	"github.com/choueric/goutils"
)

var configFile string

func init() {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		clog.Fatal("$HOME is empty")
	}
	configFile = path.Join(homeDir, ".regGen/license")
}

func Load(filepath string) (string, error) {
	if filepath == "" {
		exist, err := goutils.IsFileExist(configFile)
		if err != nil {
			return "", err
		}
		if exist {
			filepath = configFile
		}
	} else {
		exist, err := goutils.IsFileExist(filepath)
		if err != nil {
			return "", err
		}
		if !exist {
			clog.Fatal(filepath + " does not exist")
		}
	}

	if filepath == "" {
		return "", nil
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	} else {
		return string(content), nil
	}
}
