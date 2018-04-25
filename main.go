package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/choueric/clog"
	"github.com/choueric/regGen/dbg"
	"github.com/choueric/regGen/format"
	"github.com/choueric/regGen/licenseload"
	"github.com/choueric/regGen/regjar"
)

var (
	input       string
	formatArg   string
	licenseFile string
	version     = "0.0.4"
	BUILD_INFO  = ""
)

func init() {
	flag.BoolVar(&dbg.True, "d", false, "enable debug")
	flag.StringVar(&input, "i", "input.regs", "input file.")
	flag.StringVar(&formatArg, "f", "cmacro", "output format type.")
	flag.StringVar(&licenseFile, "l", "", "specify license file.")

	defUsage := flag.Usage
	flag.Usage = func() {
		fmt.Println("version:", version, BUILD_INFO)
		defUsage()
	}
	flag.Parse()

	if dbg.True {
		clog.SetFlags(clog.Lshortfile | clog.Lcolor)
		clog.Println(input)
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

	jar, err := regjar.New(input)
	if err != nil {
		clog.Fatal(err)
	}

	if dbg.True {
		fmt.Println("----------------- format output ---------------")
	}
	fmtter.FormatLicense(os.Stdout, license)
	fmtter.FormatRegJar(os.Stdout, jar)
}
