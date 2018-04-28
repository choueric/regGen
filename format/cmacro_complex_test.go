package format_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/choueric/clog"
	"github.com/choueric/goutils"
	"github.com/choueric/regGen/dbg"
	"github.com/choueric/regGen/format"
	"github.com/choueric/regGen/regjar"
)

const cmacroCmplx = `#pragma once

#ifndef BIT
#define BIT(x) (1 << (x))
#endif

// ONLY for _8bit-width_ register
#define MASK(a, b) (((uint8_t)-1 >> (7-(b))) & ~((1U<<(a))-1))

// Registers of complexChip

////////////////////////////////////////////////////////
// PLL
#define REG_PLL_BASE_ADDR 0x30000000

#define REG_PLL_CONTROL 0x0 // 0
	#define REG_PLL_BYPASS_BIT BIT(1)
	#define REG_PLL_BYPASS_POS 1
	#define REG_PLL_BYPASS_VAL(rv) (((rv) & BIT(1)) >> 1)
		#define REG_PLL_BYPASS_ENABLE	1	// 0b1	0x1
		#define REG_PLL_BYPASS_DISABLE	0	// 0b0	0x0
	#define REG_PLL_FREE_RUN_STR 6
	#define REG_PLL_FREE_RUN_END 7
	#define REG_PLL_FREE_RUN_MSK MASK(6, 7)
	#define REG_PLL_FREE_RUN_VAL(rv) (((rv) & REG_FREE_RUN_MSK) >> 6)
	#define REG_PLL_FREE_RUN_SFT(v) (((v) & MASK(0, 1)) << 6)

#define REG_PLL_2 0x2 // 2
	#define REG_PLL_TYPE_STR 4
	#define REG_PLL_TYPE_END 5
	#define REG_PLL_TYPE_MSK MASK(4, 5)
	#define REG_PLL_TYPE_VAL(rv) (((rv) & REG_TYPE_MSK) >> 4)
	#define REG_PLL_TYPE_SFT(v) (((v) & MASK(0, 1)) << 4)
		#define REG_PLL_TYPE_CLIENT	0	// 0b0	0x0
		#define REG_PLL_TYPE_SERVER	1	// 0b1	0x1
		#define REG_PLL_TYPE_ROUTE	2	// 0b10	0x2
		#define REG_PLL_TYPE_PEER	3	// 0b11	0x3

////////////////////////////////////////////////////////
// I2C
#define REG_I2C_BASE_ADDR 0x40000000

#define REG_I2C_CLOCK 0x0 // 0
	#define REG_I2C_ENABLE_BIT BIT(0)
	#define REG_I2C_ENABLE_POS 0
	#define REG_I2C_ENABLE_VAL(rv) ((rv) & BIT(0))

#define REG_I2C_DATA 0x1 // 1
	#define REG_I2C_DATA_STR 0
	#define REG_I2C_DATA_END 7
	#define REG_I2C_DATA_MSK MASK(0, 7)
	#define REG_I2C_DATA_VAL(rv) (rv)
	#define REG_I2C_DATA_SFT(v) (v)

#define REG_I2C_CONTROL 0x2 // 2
	#define REG_I2C_START_BIT BIT(0)
	#define REG_I2C_START_POS 0
	#define REG_I2C_START_VAL(rv) ((rv) & BIT(0))
		#define REG_I2C_START_START	1	// 0b1	0x1
		#define REG_I2C_START_STOP	0	// 0b0	0x0
	#define REG_I2C_INTERUPT_BIT BIT(1)
	#define REG_I2C_INTERUPT_POS 1
	#define REG_I2C_INTERUPT_VAL(rv) (((rv) & BIT(1)) >> 1)
`

func TestCmacroComplex(t *testing.T) {
	fmtter, err := format.New("cmacro")
	if err != nil {
		clog.Fatal(err)
	}

	jar, err := regjar.New("../testdata/complex.regs")
	if err != nil {
		clog.Fatal(err)
	}

	var result bytes.Buffer
	fmtter.FormatRegJar(&result, jar)
	if strings.Compare(result.String(), cmacroCmplx) != 0 {
		goutils.PrintStringMismatch(result.Bytes(), []byte(cmacroCmplx), dbg.True)
		t.Fatal("complex cmarco output does not match!")
	}
}
