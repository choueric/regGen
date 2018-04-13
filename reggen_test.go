package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/choueric/clog"
)

const (
	// Trimmed String
	trimedStr = `[  CHIP ] <chip>:si5324
[  REG  ] <REG>[Control]: 0
[ FIELD ] BYPASS_REG: 1 ()
[ FIELD ] FREE_RUN: 7 - 6 ()
[ FIELD ] ck_prior1 : 0-1 ()
[  REG  ] <REG>: 0x10
[ FIELD ] BWSEL_REG: 4-7 ()
[  REG  ] <REG>[field_vals]: 0x11
[ FIELD ] fos: 5-6 (0:fos_0, 3:fos_3)
[ FIELD ] VALTIME: 4 -  3 (0b00:    0ms, 0b01:    1ms, 0b10: 2ms, 0b11: 3ms)
[ FIELD ] lockt: 0-2 (0x0: 0t, 0xa: 5t)
[  REG  ] <reg>: 0x12
[ FIELD ] single: 2-3 (2: two)
`
	// Parsed String
	parsedStr = `CHIP: "si5324"
"Control", 0x0
    BYPASS_REG: [1:1] ()
    FREE_RUN: [6:7] ()
    ck_prior1: [0:1] ()
"16", 0x10
    BWSEL_REG: [4:7] ()
"field_vals", 0x11
    fos: [5:6] (0: fos_0, 3: fos_3, )
    VALTIME: [3:4] (0: 0ms, 1: 1ms, 2: 2ms, 3: 3ms, )
    lockt: [0:2] (0: 0t, 10: 5t, )
"18", 0x12
    single: [2:3] (2: two, )
`
	// output C format string
	formatCStr = `#pragma once

#ifndef BIT
#define BIT(x) (1 << (x))
#endif

// ONLY for _8bit-width_ register
#define MASK(a, b) (((uint8_t)-1 >> (7-(b))) & ~((1U<<(a))-1))

// Registers of si5324

#define REG_CONTROL 0x0 // 0
	#define REG_BYPASS_REG_BIT BIT(1)
	#define REG_BYPASS_REG_VAL(rv) (((rv) & BIT(1)) >> 1)
	#define REG_BYPASS_REG_POS 1
	#define REG_FREE_RUN_STR 6
	#define REG_FREE_RUN_END 7
	#define REG_FREE_RUN_MSK MASK(6, 7)
	#define REG_FREE_RUN_VAL(rv) (((rv) & REG_FREE_RUN_MSK) >> 6)
	#define REG_FREE_RUN_SFT(v) (((v) & REG_FREE_RUN_MSK) << 6)
	#define REG_CK_PRIOR1_STR 0
	#define REG_CK_PRIOR1_END 1
	#define REG_CK_PRIOR1_MSK MASK(0, 1)
	#define REG_CK_PRIOR1_VAL(rv) ((rv) & REG_CK_PRIOR1_MSK)
	#define REG_CK_PRIOR1_SFT(v) ((v) & REG_CK_PRIOR1_MSK)

#define REG_16 0x10 // 16
	#define REG_BWSEL_REG_STR 4
	#define REG_BWSEL_REG_END 7
	#define REG_BWSEL_REG_MSK MASK(4, 7)
	#define REG_BWSEL_REG_VAL(rv) (((rv) & REG_BWSEL_REG_MSK) >> 4)
	#define REG_BWSEL_REG_SFT(v) (((v) & REG_BWSEL_REG_MSK) << 4)

#define REG_FIELD_VALS 0x11 // 17
	#define REG_FOS_STR 5
	#define REG_FOS_END 6
	#define REG_FOS_MSK MASK(5, 6)
	#define REG_FOS_VAL(rv) (((rv) & REG_FOS_MSK) >> 5)
	#define REG_FOS_SFT(v) (((v) & REG_FOS_MSK) << 5)
		#define REG_FOS_FOS_0	0	// 0b0	0x0
		#define REG_FOS_FOS_3	3	// 0b11	0x3
	#define REG_VALTIME_STR 3
	#define REG_VALTIME_END 4
	#define REG_VALTIME_MSK MASK(3, 4)
	#define REG_VALTIME_VAL(rv) (((rv) & REG_VALTIME_MSK) >> 3)
	#define REG_VALTIME_SFT(v) (((v) & REG_VALTIME_MSK) << 3)
		#define REG_VALTIME_0MS	0	// 0b0	0x0
		#define REG_VALTIME_1MS	1	// 0b1	0x1
		#define REG_VALTIME_2MS	2	// 0b10	0x2
		#define REG_VALTIME_3MS	3	// 0b11	0x3
	#define REG_LOCKT_STR 0
	#define REG_LOCKT_END 2
	#define REG_LOCKT_MSK MASK(0, 2)
	#define REG_LOCKT_VAL(rv) ((rv) & REG_LOCKT_MSK)
	#define REG_LOCKT_SFT(v) ((v) & REG_LOCKT_MSK)
		#define REG_LOCKT_0T	0	// 0b0		0x0
		#define REG_LOCKT_5T	10	// 0b1010	0xa

#define REG_18 0x12 // 18
	#define REG_SINGLE_STR 2
	#define REG_SINGLE_END 3
	#define REG_SINGLE_MSK MASK(2, 3)
	#define REG_SINGLE_VAL(rv) (((rv) & REG_SINGLE_MSK) >> 2)
	#define REG_SINGLE_SFT(v) (((v) & REG_SINGLE_MSK) << 2)
		#define REG_SINGLE_TWO	2	// 0b10	0x2
`
)

var (
	testSource string
	testItems  []*lineItem
	testDebug  bool
)

func init() {
	data, err := ioutil.ReadFile("./chips/test.regs")
	if err != nil {
		clog.Fatal(err)
	}
	testSource = string(data)

	flag.BoolVar(&testDebug, "d", false, "enable test debug")
	flag.Parse()

	testItems, err = trim(bufio.NewReader(strings.NewReader(testSource)))
	if err != nil {
		clog.Fatal(err)
	}
}

func print_string_mismatch(a, b []byte) {
	if len(a) != len(b) {
		fmt.Println("length doesn't match:", len(a), len(b))
		if testDebug {
			fmt.Println(string(a), "----------------------\n", string(b))
		}
		return
	}
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			fmt.Println(i, string(a[i]), string(b[i]))
			break
		}
	}
}

func Test_trim(t *testing.T) {
	// testItems is already got in init()
	var result bytes.Buffer
	printTrimItems(&result, testItems)
	if strings.Compare(result.String(), trimedStr) != 0 {
		print_string_mismatch(result.Bytes(), []byte(trimedStr))
		t.Error("trim fail!")
	}
}

func Test_parse(t *testing.T) {
	jar, err := parse(testItems)
	if err != nil {
		t.Error(err)
	}

	var result bytes.Buffer
	fmt.Fprint(&result, jar)
	if strings.Compare(result.String(), parsedStr) != 0 {
		print_string_mismatch(result.Bytes(), []byte(parsedStr))
		t.Error("parse fail!")
	}
}

func Test_formatToC(t *testing.T) {
	jar, err := parse(testItems)
	if err != nil {
		t.Error(err)
	}

	var result bytes.Buffer
	formatToC(jar, &result)
	if strings.Compare(result.String(), formatCStr) != 0 {
		print_string_mismatch(result.Bytes(), []byte(formatCStr))
		t.Error("parse fail!")
	}
}
