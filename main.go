package main

import (
	"fmt"
	"os"
	"path"

	"github.com/choueric/clog"
	"github.com/choueric/regGen/dbg"
	"github.com/choueric/regGen/fileflag"
	"github.com/choueric/regGen/format"
	"github.com/choueric/regGen/licenseload"
	"github.com/choueric/regGen/regjar"
)

var (
	inputArg       string
	formatArg      string
	showVersionArg bool
	isFullArg      bool
	licenseFile    string
	version        = "0.0.4"
	BUILD_INFO     = ""
)

func joinHomeDir(filepath string) string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		clog.Fatal("$HOME is empty")
	}
	return path.Join(homeDir, filepath)
}

func init() {
	ffPath := joinHomeDir(".regGen/flag")
	ff := fileflag.New(ffPath)

	ff.FlagSet().BoolVar(&dbg.True, "d", false, "enable debug")
	ff.FlagSet().BoolVar(&showVersionArg, "v", false, "show version.")
	ff.FlagSet().StringVar(&inputArg, "i", "input.regs", "input file.")
	ff.FlagSet().StringVar(&formatArg, "f", "cmacro", "output format type.")
	ff.FlagSet().StringVar(&licenseFile, "l", "", "specify license file.")
	ff.FlagSet().BoolVar(&isFullArg, "full", false, "full format output")
	if err := ff.Parse(); err != nil {
		clog.Fatal(err)
	}

	if showVersionArg {
		fmt.Println("version:", version, BUILD_INFO)
		os.Exit(0)
	}

	if dbg.True {
		clog.SetFlags(clog.Lshortfile | clog.Lcolor)
		clog.Println(inputArg)
	}
}

func main() {
	fmtter, err := format.New(formatArg)
	if err != nil {
		clog.Fatal(err)
	}

	license, err := licenseload.Load(licenseFile)
	if err != nil {
		clog.Fatal(err)
	}

	jar, err := regjar.New(inputArg)
	if err != nil {
		clog.Fatal(err)
	}

	if dbg.True {
		fmt.Println("----------------- format output ---------------")
	}

	fmtter.FormatLicense(os.Stdout, license)
	fmtter.FormatRegJar(os.Stdout, jar, isFullArg)
}
