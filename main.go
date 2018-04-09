package main

import (
	"flag"
	"os"

	"github.com/choueric/clog"
)

var (
	input  string
	debug  bool
	format string
)

func main() {
	flag.BoolVar(&debug, "d", false, "enable debug")
	flag.StringVar(&input, "i", "input.regs", "input file.")
	flag.StringVar(&format, "f", "c", "output format type. [c]")
	flag.Parse()
	if len(os.Args) == 2 && (os.Args[1] == "help" || os.Args[1] == "-h") {
		flag.Usage()
		return
	}

	if debug {
		clog.Println(input)
	}

	var regs regMap
	err := regs.Load(input)
	if err != nil {
		clog.Fatal(err)
	}

	regs.Output(os.Stdout, format)
}
