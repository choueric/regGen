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
	version    = "0.0.1"
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

	var regs regMap
	err := regs.Load(input)
	if err != nil {
		clog.Fatal(err)
	}

	regs.Output(os.Stdout, format)
}
