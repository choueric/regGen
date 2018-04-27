package format

import (
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"

	"github.com/choueric/clog"
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
	if len(jar.Modules) == 0 {
		clog.Fatal("Empty modules in jar\n")
	}

	tw := tabwriter.NewWriter(w, 0, 4, 1, '\t', 0)
	fmtter.formatBanner(tw, jar)
	if len(jar.Modules) == 1 && jar.Modules[0].Name == "default" {
		for _, r := range jar.Modules[0].Regs {
			prefix := "#define REG"
			fmtter.formatReg(w, r, prefix)
			tw.Flush()
		}
	} else {
		for _, mod := range jar.Modules {
			prefix := "#define REG"
			fmtter.formatModule(tw, mod, prefix)
			prefix = fmt.Sprintf("%s_%s", prefix, mod.Name)
			for _, r := range mod.Regs {
				fmtter.formatReg(w, r, prefix)
			}
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

func (fmtter *cmacroFormat) formatModule(w io.Writer, m *regjar.Module, prefix string) {
	fmt.Fprintf(w, "\n")
	fmt.Fprintln(w, "////////////////////////////////////////////////////////")
	fmt.Fprintf(w, "// %s\n", m.Name)
	fmt.Fprintf(w, "%s_%s_BASE_ADDR 0x%08x\n", prefix, m.Name, m.BaseAddr)
}

func (fmtter *cmacroFormat) formatReg(w io.Writer, r *regjar.Reg, prefix string) {
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, prefix+"_%s %#x // %d\n", r.Name, r.Offset, r.Offset)
	for _, f := range r.Fields {
		fmtter.formatField(w, f, "\t"+prefix)
	}
}

func (fmtter *cmacroFormat) formatField(w io.Writer, f *regjar.Field, prefix string) {
	if f.Start == f.End {
		fmtter.formatBitField(w, f, prefix)
	} else {
		fmtter.formatRangeField(w, f, prefix)
	}
	fmtter.formatEnums(w, f, "\t"+prefix)
}

func (fmtter *cmacroFormat) formatBitField(w io.Writer, f *regjar.Field, prefix string) {
	p := f.Start
	fmt.Fprintf(w, "%s_%s_BIT BIT(%d)\n", prefix, f.Name, p)
	fmt.Fprintf(w, "%s_%s_POS %d\n", prefix, f.Name, p)
	if p == 0 {
		fmt.Fprintf(w, "%s_%s_VAL(rv) ((rv) & BIT(%d))\n", prefix, f.Name, p)
	} else {
		fmt.Fprintf(w, "%s_%s_VAL(rv) (((rv) & BIT(%d)) >> %d)\n", prefix, f.Name, p, p)
	}
}

func (fmtter *cmacroFormat) formatRangeField(w io.Writer, f *regjar.Field, prefix string) {
	// start & end
	fmt.Fprintf(w, "%s_%s_STR %d\n", prefix, f.Name, f.Start)
	fmt.Fprintf(w, "%s_%s_END %d\n", prefix, f.Name, f.End)

	// mask
	fmt.Fprintf(w, "%s_%s_MSK MASK(%d, %d)\n", prefix, f.Name, f.Start, f.End)

	// val
	if f.Start == 0 && f.End == 7 {
		fmt.Fprintf(w, "%s_%s_VAL(rv) (rv)\n", prefix, f.Name)
	} else if f.Start == 0 {
		fmt.Fprintf(w, "%s_%s_VAL(rv) ((rv) & REG_%s_MSK)\n", prefix, f.Name, f.Name)
	} else {
		fmt.Fprintf(w, "%s_%s_VAL(rv) (((rv) & REG_%s_MSK) >> %d)\n",
			prefix, f.Name, f.Name, f.Start)
	}

	// shift
	if f.Start == 0 && f.End == 7 {
		fmt.Fprintf(w, "%s_%s_SFT(v) (v)\n", prefix, f.Name)
	} else if f.Start == 0 {
		fmt.Fprintf(w, "%s_%s_SFT(v) ((v) & MASK(0, %d))\n", prefix, f.Name, f.End)
	} else {
		fmt.Fprintf(w, "%s_%s_SFT(v) (((v) & MASK(0, %d)) << %d)\n",
			prefix, f.Name, f.End-f.Start, f.Start)
	}
}

func (fmtter *cmacroFormat) formatEnums(w io.Writer, f *regjar.Field, prefix string) {
	if len(f.EnumVals) != 0 {
		for i, v := range f.EnumVals {
			bstr := strconv.FormatUint(uint64(v), 2)
			fmt.Fprintf(w, "%s_%s_%s\t%d\t// 0b%s\t%#x\n",
				prefix, f.Name, f.EnumNames[i], v, bstr, v)
		}
	}
}
