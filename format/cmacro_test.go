package format_test

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/choueric/clog"
	"github.com/choueric/goutils"
	"github.com/choueric/regGen/dbg"
	"github.com/choueric/regGen/format"
	"github.com/choueric/regGen/regjar"
)

const formatCStr = `#pragma once

#ifndef BIT
#define BIT(x) (1 << (x))
#endif

// ONLY for _8bit-width_ register
#define MASK(a, b) (((uint8_t)-1 >> (7-(b))) & ~((1U<<(a))-1))

// Registers of si5324

#define REG_CONTROL 0x0 // 0
	#define REG_BYPASS_REG_BIT BIT(1)
	#define REG_BYPASS_REG_POS 1
	#define REG_BYPASS_REG_VAL(rv) (((rv) & BIT(1)) >> 1)
	#define REG_FREE_RUN_STR 6
	#define REG_FREE_RUN_END 7
	#define REG_FREE_RUN_MSK MASK(6, 7)
	#define REG_FREE_RUN_VAL(rv) (((rv) & REG_FREE_RUN_MSK) >> 6)
	#define REG_FREE_RUN_SFT(v) (((v) & MASK(0, 1)) << 6)
	#define REG_CK_PRIOR1_STR 0
	#define REG_CK_PRIOR1_END 1
	#define REG_CK_PRIOR1_MSK MASK(0, 1)
	#define REG_CK_PRIOR1_VAL(rv) ((rv) & REG_CK_PRIOR1_MSK)
	#define REG_CK_PRIOR1_SFT(v) ((v) & MASK(0, 1))

#define REG_16 0x10 // 16
	#define REG_BWSEL_REG_STR 4
	#define REG_BWSEL_REG_END 7
	#define REG_BWSEL_REG_MSK MASK(4, 7)
	#define REG_BWSEL_REG_VAL(rv) (((rv) & REG_BWSEL_REG_MSK) >> 4)
	#define REG_BWSEL_REG_SFT(v) (((v) & MASK(0, 3)) << 4)

#define REG_FIELD_VALS 0x11 // 17
	#define REG_FOS_STR 5
	#define REG_FOS_END 6
	#define REG_FOS_MSK MASK(5, 6)
	#define REG_FOS_VAL(rv) (((rv) & REG_FOS_MSK) >> 5)
	#define REG_FOS_SFT(v) (((v) & MASK(0, 1)) << 5)
		#define REG_FOS_FOS_0	0	// 0b0	0x0
		#define REG_FOS_FOS_3	3	// 0b11	0x3
	#define REG_VALTIME_STR 3
	#define REG_VALTIME_END 4
	#define REG_VALTIME_MSK MASK(3, 4)
	#define REG_VALTIME_VAL(rv) (((rv) & REG_VALTIME_MSK) >> 3)
	#define REG_VALTIME_SFT(v) (((v) & MASK(0, 1)) << 3)
		#define REG_VALTIME_0MS	0	// 0b0	0x0
		#define REG_VALTIME_1MS	1	// 0b1	0x1
		#define REG_VALTIME_2MS	2	// 0b10	0x2
		#define REG_VALTIME_3MS	3	// 0b11	0x3
	#define REG_LOCKT_STR 0
	#define REG_LOCKT_END 2
	#define REG_LOCKT_MSK MASK(0, 2)
	#define REG_LOCKT_VAL(rv) ((rv) & REG_LOCKT_MSK)
	#define REG_LOCKT_SFT(v) ((v) & MASK(0, 2))
		#define REG_LOCKT_0T	0	// 0b0		0x0
		#define REG_LOCKT_5T	10	// 0b1010	0xa

#define REG_18 0x12 // 18
	#define REG_SINGLE_STR 2
	#define REG_SINGLE_END 3
	#define REG_SINGLE_MSK MASK(2, 3)
	#define REG_SINGLE_VAL(rv) (((rv) & REG_SINGLE_MSK) >> 2)
	#define REG_SINGLE_SFT(v) (((v) & MASK(0, 1)) << 2)
		#define REG_SINGLE_TWO	2	// 0b10	0x2
`

func init() {
	flag.BoolVar(&dbg.True, "d", false, "enable test debug")
	flag.Parse()
	if dbg.True {
		clog.SetFlags(clog.Lshortfile | clog.Lcolor)
	}
}

func TestCMacrosFormat(t *testing.T) {
	fmtter, err := format.New("cmacro")
	if err != nil {
		clog.Fatal(err)
	}

	jar, err := regjar.New("../chips/test.regs")
	if err != nil {
		clog.Fatal(err)
	}

	var result bytes.Buffer
	fmtter.FormatRegJar(&result, jar)
	if strings.Compare(result.String(), formatCStr) != 0 {
		goutils.PrintStringMismatch(result.Bytes(), []byte(formatCStr), dbg.True)
		t.Error("parse fail!")
	}
	fmt.Println("ok")
}
