package main

import (
	"fmt"
	"io"
	"strings"
)

type formatFunc func(rm *regMap, w io.Writer)

var outputFormat = map[string]formatFunc{
	"c": formatToC,
}

const cHeader = `#pragma once

#ifndef BIT
#define BIT(x) (1 << (x))
#endif

// ONLY for _8bit-width_ register
#define MASK(a, b) (((uint8_t)-1 >> (7-(b))) & ~((1U<<(a))-1))
`

/*
#define REG_0 0x0
	#define REG_FREE_RUN_BIT BIT(1)
	#define REG_CKOUT_ALWARYS_ON_BIT BIT(5)
	#define REG_CKSEL_MSK MASK(6, 7)
*/
func formatToC(rm *regMap, w io.Writer) {
	fmt.Fprintf(w, cHeader)
	if rm.chip != "" {
		fmt.Fprintf(w, "\n// Registers of %s\n", rm.chip)
	}
	for _, r := range rm.regs {
		fmt.Fprintf(w, "\n#define REG_%s %#x // %d\n", strings.ToUpper(r.name),
			r.offset, r.offset)
		for _, f := range r.fields {
			if f.start == f.end {
				fmt.Fprintf(w, "\t#define REG_%s_BIT BIT(%d)\n",
					strings.ToUpper(f.name), f.start)
			} else {
				fmt.Fprintf(w, "\t#define REG_%s_MSK MASK(%d, %d)\n",
					strings.ToUpper(f.name), f.start, f.end)
			}
		}
	}
}
