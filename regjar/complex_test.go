package regjar

import (
	"bufio"
	"bytes"
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
	cmplxTrimStr = `[  CHIP ] <chip>:complexChip
[  MOD  ] <module>[pll]: 0x30000000
[  REG  ] <REG>[Control]: 0
[ FIELD ] BYPASS: 1 (1: enable, 0: disable)
[ FIELD ] FREE_RUN: 7 - 6
[  REG  ] <REG>: 0x2
[ FIELD ] type: 4-5 (0b00: client, 0x01: server, 2: route, 3: peer)
[  MOD  ] <module>[i2c]: 0x40000000
[  REG  ] <REG>[clock]: 0
[ FIELD ] enable: 0
[  REG  ] <REG>[data]: 1
[ FIELD ] data: 0-7
[  REG  ] <REG>[control]:2
[ FIELD ] start: 0 (1: start, 0: stop)
[ FIELD ] interupt: 1
`
	// Parsed String
	cmplxParseStr = `CHIP: "complexChip"
WIDTH: 8
MODULE: "PLL"
"CONTROL", 0x0
    BYPASS: [1:1] (1: ENABLE, 0: DISABLE, )
    FREE_RUN: [6:7] ()
"2", 0x2
    TYPE: [4:5] (0: CLIENT, 1: SERVER, 2: ROUTE, 3: PEER, )
MODULE: "I2C"
"CLOCK", 0x0
    ENABLE: [0:0] ()
"DATA", 0x1
    DATA: [0:7] ()
"CONTROL", 0x2
    START: [0:0] (1: START, 0: STOP, )
    INTERUPT: [1:1] ()
`
)

var (
	cmplxSource string
	cmplxItems  tagItemSlice
)

func init() {
	data, err := ioutil.ReadFile("../testdata/complex.regs")
	if err != nil {
		clog.Fatal(err)
	}
	cmplxSource = string(data)

	cmplxItems, err = trim(bufio.NewReader(strings.NewReader(cmplxSource)))
	if err != nil {
		clog.Fatal(err)
	}
}

func TestComplexTrim(t *testing.T) {
	// cmplxItems is already got in init()
	result := cmplxItems.String()
	if strings.Compare(result, cmplxTrimStr) != 0 {
		goutils.PrintStringMismatch([]byte(result), []byte(cmplxTrimStr), dbg.True)
		t.Error("trim fail!")
	}
}

func TestComplexParse(t *testing.T) {
	jar, err := parse(cmplxItems)
	if err != nil {
		t.Error(err)
	}

	var result bytes.Buffer
	fmt.Fprint(&result, jar)
	if strings.Compare(result.String(), cmplxParseStr) != 0 {
		goutils.PrintStringMismatch(result.Bytes(), []byte(cmplxParseStr), dbg.True)
		t.Error("parse fail!")
	}
}
