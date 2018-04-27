package regjar

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/choueric/clog"
	"github.com/choueric/goutils"
	"github.com/choueric/regGen/dbg"
)

const (
	// Trimmed String
	trimedStr = `[  CHIP ] <chip>:si5324
[  REG  ] <REG>[Control]: 0
[ FIELD ] BYPASS_REG: 1
[ FIELD ] FREE_RUN: 7 - 6
[ FIELD ] ck_prior1 : 0-1
[  REG  ] <REG>: 0x10
[ FIELD ] BWSEL_REG: 4-7
[  REG  ] <REG>[field_vals]: 0x11
[ FIELD ] fos: 5-6 (0:fos_0, 3:fos_3)
[ FIELD ] VALTIME: 4 -  3 (0b00:    0ms, 0b01:    1ms, 0b10: 2ms, 0b11: 3ms)
[ FIELD ] lockt: 0-2 (0x0: 0t, 0xa: 5t)
[  REG  ] <reg>: 0x12
[ FIELD ] single: 2-3 (2: two)
`
	// Parsed String
	parsedStr = `CHIP: "si5324"
WIDTH: 8
MODULE: "default"
"CONTROL", 0x0
    BYPASS_REG: [1:1] ()
    FREE_RUN: [6:7] ()
    CK_PRIOR1: [0:1] ()
"16", 0x10
    BWSEL_REG: [4:7] ()
"FIELD_VALS", 0x11
    FOS: [5:6] (0: FOS_0, 3: FOS_3, )
    VALTIME: [3:4] (0: 0MS, 1: 1MS, 2: 2MS, 3: 3MS, )
    LOCKT: [0:2] (0: 0T, 10: 5T, )
"18", 0x12
    SINGLE: [2:3] (2: TWO, )
`
)

var (
	testSource string
	testItems  tagItemSlice
)

func init() {
	data, err := ioutil.ReadFile("../testdata/test.regs")
	if err != nil {
		clog.Fatal(err)
	}
	testSource = string(data)

	flag.BoolVar(&dbg.True, "d", false, "enable test debug")
	flag.Parse()
	if dbg.True {
		clog.SetFlags(clog.Lshortfile | clog.Lcolor)
	}

	testItems, err = trim(bufio.NewReader(strings.NewReader(testSource)))
	if err != nil {
		clog.Fatal(err)
	}
}

func TestTrim(t *testing.T) {
	// testItems is already got in init()
	result := testItems.String()
	if strings.Compare(result, trimedStr) != 0 {
		goutils.PrintStringMismatch([]byte(result), []byte(trimedStr), dbg.True)
		t.Error("trim fail!")
	}
}

func TestParse(t *testing.T) {
	jar, err := parse(testItems)
	if err != nil {
		t.Error(err)
	}

	var result bytes.Buffer
	fmt.Fprint(&result, jar)
	if strings.Compare(result.String(), parsedStr) != 0 {
		goutils.PrintStringMismatch(result.Bytes(), []byte(parsedStr), dbg.True)
		t.Error("parse fail!")
	}
}
