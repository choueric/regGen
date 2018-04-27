package fileflag

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	var (
		debug       bool
		input       string
		formatArg   string
		licenseFile string
	)
	os.Args = []string{"fileflag.test"} // clear cmdline argH
	ff := New("../testdata/fileflag_test.conf")
	ff.FlagSet().BoolVar(&debug, "d", false, "enable debug")
	ff.FlagSet().StringVar(&input, "i", "input.regs", "input file.")
	ff.FlagSet().StringVar(&formatArg, "f", "cmacro", "output format type.")
	ff.FlagSet().StringVar(&licenseFile, "l", "", "specify license file.")

	if err := ff.Parse(); err != nil {
		t.Fatal(err)
	}

	if debug != true {
		t.Fatal("debug dosen't match.")
	}
	if input != "../chips/simple.regs" {
		t.Fatal("input dosen't match.")
	}
	if formatArg != "cmacro" {
		t.Fatal("output format dosen't match.")
	}
	if licenseFile != "../LICENSE" {
		t.Fatal("license file dosen't match.")
	}
}

func TestParseCmdArgs(t *testing.T) {
	var (
		debug       bool
		input       string
		formatArg   string
		licenseFile string
	)
	ff := New("")
	ff.FlagSet().BoolVar(&debug, "d", false, "enable debug")
	ff.FlagSet().StringVar(&input, "i", "input.regs", "input file.")
	ff.FlagSet().StringVar(&formatArg, "f", "cmacro", "output format type.")
	ff.FlagSet().StringVar(&licenseFile, "l", "", "specify license file.")

	os.Args = []string{"test", "-d", "-i", "../chips/simple.regs", "-l", "../LICENSE"}
	if err := ff.Parse(); err != nil {
		t.Fatal(err)
	}

	if debug != true {
		t.Fatal("debug dosen't match.")
	}
	if input != "../chips/simple.regs" {
		t.Fatal("input dosen't match.")
	}
	if formatArg != "cmacro" {
		t.Fatal("output format dosen't match.")
	}
	if licenseFile != "../LICENSE" {
		t.Fatal("license file dosen't match.")
	}
}
