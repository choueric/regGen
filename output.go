package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
)

type formatFunc func(jar *regJar, w io.Writer)

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

func cfmtOutputBitField(w io.Writer, f *field, n string) {
	p := f.start
	fmt.Fprintf(w, "\t#define REG_%s_BIT BIT(%d)\n", n, p)
	if p == 0 {
		fmt.Fprintf(w, "\t#define REG_%s_VAL(rv) ((rv) & BIT(%d))\n", n, p)
	} else {
		fmt.Fprintf(w, "\t#define REG_%s_VAL(rv) (((rv) & BIT(%d)) >> %d)\n", n, p, p)
	}
	fmt.Fprintf(w, "\t#define REG_%s_POS %d\n", n, p)
}

func cfmtOutputMaskField(w io.Writer, f *field, n string) {
	// start & end
	fmt.Fprintf(w, "\t#define REG_%s_STR %d\n", n, f.start)
	fmt.Fprintf(w, "\t#define REG_%s_END %d\n", n, f.end)

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
}

func cfmtOputputValues(w io.Writer, f *field, n string) {
	// values
	if len(f.valData) != 0 {
		for i, v := range f.valData {
			bstr := strconv.FormatUint(uint64(v), 2)
			fmt.Fprintf(w, "\t\t#define REG_%s_%s\t%d\t// 0b%s\t%#x\n",
				n, strings.ToUpper(f.valName[i]), v, bstr, v)
		}
	}
}

func formatToC(jar *regJar, ow io.Writer) {
	w := tabwriter.NewWriter(ow, 0, 4, 1, '\t', 0)
	fmt.Fprintf(w, cHeader)
	if jar.chip != "" {
		fmt.Fprintf(w, "\n// Registers of %s\n", jar.chip)
	}
	for _, r := range jar.regs {
		fmt.Fprintf(w, "\n#define REG_%s %#x // %d\n", strings.ToUpper(r.name),
			r.offset, r.offset)
		for _, f := range r.fields {
			name := strings.ToUpper(f.name)
			if f.start == f.end {
				cfmtOutputBitField(w, f, name)
			} else {
				cfmtOutputMaskField(w, f, name)
			}
			cfmtOputputValues(w, f, name)
		}
	}
	w.Flush()
}
