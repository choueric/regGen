package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/choueric/clog"
)

var (
	input      string
	debug      bool
	format     string
	version    = "0.0.3"
	BUILD_INFO = ""
)

func main() {
	flag.BoolVar(&debug, "d", false, "enable debug")
	flag.StringVar(&input, "i", "input.regs", "input file.")
	flag.StringVar(&format, "f", "c", "output format type. [c]")

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

	jar, err := newRegJar(input)
	if err != nil {
		clog.Fatal(err)
	}

	if debug {
		fmt.Println("----------------- format output ---------------")
	}
	fmtFunc(jar, os.Stdout)
}
