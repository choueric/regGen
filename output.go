package main

import (
	"fmt"
	"io"
	"strconv"
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

func cfmtOutputMaskField(w io.Writer, f *field, n string) {
	// mask
	fmt.Fprintf(w, "\t#define REG_%s_MSK MASK(%d, %d)\n", n, f.start, f.end)

	// val
	if f.start == 0 && f.end == 7 {
		fmt.Fprintf(w, "\t#define REG_%s_VAL(rv) (rv)\n", n)
	} else if f.start == 0 {
		fmt.Fprintf(w, "\t#define REG_%s_VAL(rv) ((rv) & REG_%s_MSK)\n", n, n)
	} else {
		fmt.Fprintf(w, "\t#define REG_%s_VAL(rv) (((rv) & REG_%s_MSK) >> %d)\n",
			n, n, f.start)
	}

	// shift
	if f.start == 0 && f.end == 7 {
		fmt.Fprintf(w, "\t#define REG_%s_SFT(v) (v)\n", n)
	} else if f.start == 0 {
		fmt.Fprintf(w, "\t#define REG_%s_SFT(v) ((v) & REG_%s_MSK)\n", n, n)
	} else {
		fmt.Fprintf(w, "\t#define REG_%s_SFT(v) (((v) & REG_%s_MSK) << %d)\n",
			n, n, f.start)
	}

	// values
	if len(f.valData) != 0 {
		for i, v := range f.valData {
			bstr := strconv.FormatUint(uint64(v), 2)
			fmt.Fprintf(w, "\t\t#define REG_%s_%s \t%d \t// 0b%s\t%#x\n",
				n, strings.ToUpper(f.valName[i]), v, bstr, v)
		}
	}
}

func formatToC(rm *regMap, w io.Writer) {
	fmt.Fprintf(w, cHeader)
	if rm.chip != "" {
		fmt.Fprintf(w, "\n// Registers of %s\n", rm.chip)
	}
	for _, r := range rm.regs {
		fmt.Fprintf(w, "\n#define REG_%s %#x // %d\n", strings.ToUpper(r.name),
			r.offset, r.offset)
		for _, f := range r.fields {
			name := strings.ToUpper(f.name)
			if f.start == f.end {
				fmt.Fprintf(w, "\t#define REG_%s_BIT BIT(%d)\n", name, f.start)
			} else {
				cfmtOutputMaskField(w, f, name)
			}
		}
	}
}
