package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/choueric/clog"
)

var (
	input             string
	debug             bool
	format            string
	licenseFile       string
	licenseConfigFile string
	version           = "0.0.4"
	BUILD_INFO        = ""
)

func init() {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		clog.Fatal("$HOME is empty")
	}
	licenseConfigFile = path.Join(homeDir, ".regGen/license")
}

func loadLicense(filepath string) (string, error) {
	if filepath == "" {
		exist, err := isFileExist(licenseConfigFile)
		if err != nil {
			return "", err
		}
		if exist {
			filepath = licenseConfigFile
		}
	} else {
		exist, err := isFileExist(filepath)
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

func isFileExist(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func main() {
	flag.BoolVar(&debug, "d", false, "enable debug")
	flag.StringVar(&input, "i", "input.regs", "input file.")
	flag.StringVar(&format, "f", "c", "output format type. [c]")
	flag.StringVar(&licenseFile, "l", "", "specify license file.")

	defUsage := flag.Usage
	flag.Usage = func() {
		fmt.Println("version:", version, BUILD_INFO)
		defUsage()
	}
	flag.Parse()

	if debug {
		clog.SetFlags(clog.Lshortfile | clog.Lcolor)
		clog.Println(input)
	}

	fmtFunc, ok := outputFormat[format]
	if !ok {
		clog.Fatal("Invalid format: " + format)
	}

	license, err := loadLicense(licenseFile)
	if err != nil {
		clog.Fatal(err)
	}

	jar, err := newRegJar(input)
	if err != nil {
		clog.Fatal(err)
	}

	if debug {
		fmt.Println("----------------- format output ---------------")
	}
	cfmtOutputLicense(os.Stdout, license)
	fmtFunc(jar, os.Stdout)
}
