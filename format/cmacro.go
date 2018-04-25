package format

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/choueric/goutils"
	"github.com/choueric/regGen/regjar"
)

type cmacroFormat int

func (fmtter *cmacroFormat) FormatLicense(w io.Writer, license string) {
	if license == "" {
		return
	}
	goutils.PrefixStringPerLine(w, license, "// ")
	fmt.Fprintf(w, "\n")
}

func (fmtter *cmacroFormat) FormatRegJar(w io.Writer, jar *regjar.Jar) {
	tw := tabwriter.NewWriter(w, 0, 4, 1, '\t', 0)
	fmtter.formatBanner(tw, jar)
	for _, r := range jar.Regs {
		fmtter.formatReg(tw, r)
		for _, f := range r.Fields {
			name := strings.ToUpper(f.Name)
			if f.Start == f.End {
				fmtter.formatBitField(tw, f, name)
			} else {
				fmtter.formatRangeField(tw, f, name)
			}
			fmtter.formatEnums(tw, f, name)
			tw.Flush()
		}
	}
	tw.Flush()
}

func (fmtter *cmacroFormat) formatBanner(w io.Writer, jar *regjar.Jar) {
	const cHeader = `#pragma once

#ifndef BIT
#define BIT(x) (1 << (x))
#endif

// ONLY for _8bit-width_ register
#define MASK(a, b) (((uint8_t)-1 >> (7-(b))) & ~((1U<<(a))-1))
`

	fmt.Fprintf(w, cHeader)
	if jar.Chip != "" {
		fmt.Fprintf(w, "\n// Registers of %s\n", jar.Chip)
	}
}

func (fmtter *cmacroFormat) formatReg(w io.Writer, r *regjar.Reg) {
	n := strings.ToUpper(r.Name)
	fmt.Fprintf(w, "\n#define REG_%s %#x // %d\n", n, r.Offset, r.Offset)
}

func (fmtter *cmacroFormat) formatBitField(w io.Writer, f *regjar.Field, n string) {
	p := f.Start
	fmt.Fprintf(w, "\t#define REG_%s_BIT BIT(%d)\n", n, p)
	fmt.Fprintf(w, "\t#define REG_%s_POS %d\n", n, p)
	if p == 0 {
		fmt.Fprintf(w, "\t#define REG_%s_VAL(rv) ((rv) & BIT(%d))\n", n, p)
	} else {
		fmt.Fprintf(w, "\t#define REG_%s_VAL(rv) (((rv) & BIT(%d)) >> %d)\n", n, p, p)
	}
}

func (fmtter *cmacroFormat) formatRangeField(w io.Writer, f *regjar.Field, n string) {
	// start & end
	fmt.Fprintf(w, "\t#define REG_%s_STR %d\n", n, f.Start)
	fmt.Fprintf(w, "\t#define REG_%s_END %d\n", n, f.End)

	// mask
	fmt.Fprintf(w, "\t#define REG_%s_MSK MASK(%d, %d)\n", n, f.Start, f.End)

	// val
	if f.Start == 0 && f.End == 7 {
		fmt.Fprintf(w, "\t#define REG_%s_VAL(rv) (rv)\n", n)
	} else if f.Start == 0 {
		fmt.Fprintf(w, "\t#define REG_%s_VAL(rv) ((rv) & REG_%s_MSK)\n", n, n)
	} else {
		fmt.Fprintf(w, "\t#define REG_%s_VAL(rv) (((rv) & REG_%s_MSK) >> %d)\n",
			n, n, f.Start)
	}

	// shift
	if f.Start == 0 && f.End == 7 {
		fmt.Fprintf(w, "\t#define REG_%s_SFT(v) (v)\n", n)
	} else if f.Start == 0 {
		fmt.Fprintf(w, "\t#define REG_%s_SFT(v) ((v) & MASK(0, %d))\n", n, f.End)
	} else {
		fmt.Fprintf(w, "\t#define REG_%s_SFT(v) (((v) & MASK(0, %d)) << %d)\n",
			n, f.End-f.Start, f.Start)
	}
}

func (fmtter *cmacroFormat) formatEnums(w io.Writer, f *regjar.Field, n string) {
	// Enums
	if len(f.EnumVals) != 0 {
		for i, v := range f.EnumVals {
			bstr := strconv.FormatUint(uint64(v), 2)
			fmt.Fprintf(w, "\t\t#define REG_%s_%s\t%d\t// 0b%s\t%#x\n",
				n, strings.ToUpper(f.EnumNames[i]), v, bstr, v)
		}
	}
}
