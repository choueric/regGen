package regjar

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/choueric/clog"
	"github.com/choueric/goutils"
	"github.com/choueric/regGen/dbg"
)

type Jar struct {
	Chip  string
	Width uint32
	Regs  []*Reg
}

func (jar *Jar) String() string {
	var str bytes.Buffer
	fmt.Fprintf(&str, "CHIP: \"%s\"\n", jar.Chip)
	fmt.Fprintf(&str, "WIDTH: %d\n", jar.Width)
	for _, r := range jar.Regs {
		fmt.Fprint(&str, r)
	}
	return str.String()
}

func (jar *Jar) addRegs(v ...*Reg) {
	jar.Regs = append(jar.Regs, v...)
}

// create a new jar and setup default values
func newJar() *Jar {
	jar := new(Jar)
	jar.Width = 8
	return jar
}

func processChip(line string) (string, error) {
	strs := strings.Split(line, ":")
	if len(strs) != 2 {
		clog.Fatal("Invalid Format: [" + line + "]")
	}
	return strings.TrimSpace(strs[1]), nil
}

func processWidth(line string) (uint32, error) {
	strs := strings.Split(line, ":")
	if len(strs) != 2 {
		clog.Fatal("Invalid Format: [" + line + "]")
	}

	w, err := goutils.ParseUint(strings.TrimSpace(strs[1]), 32)
	if err != nil {
		return 0, err
	}

	if w != 8 && w != 16 && w != 32 && w != 64 {
		return 0, errors.New(fmt.Sprintf("Invalid bit width %d", w))
	}
	return uint32(w), nil
}

func trim(reader *bufio.Reader) (tagItemSlice, error) {
	items := newTagItemSlice()
	for {
		line, err := goutils.ReadLine(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		item := newTagItem(line)
		if err != nil {
			clog.Fatal(err)
		} else {
			if item != nil && item.tag != tag_comment && item.tag != tag_other {
				items.addTagItems(item)
			}
		}
	}

	if dbg.True {
		fmt.Println("----------------- after trim ---------------")
		fmt.Println(items)
	}

	return items, nil
}

func parse(items tagItemSlice) (*Jar, error) {
	var curReg *Reg
	jar := newJar()
	for _, item := range items {
		switch item.tag {
		case tag_chip:
			chip, err := processChip(item.data)
			if err != nil {
				return nil, err
			}
			jar.Chip = chip
		case tag_width:
			width, err := processWidth(item.data)
			if err != nil {
				return nil, err
			}
			jar.Width = width
		case tag_reg:
			r, err := processReg(item.data)
			if err != nil {
				return nil, err
			}
			jar.addRegs(r)
			curReg = r
		case tag_field:
			if curReg == nil {
				clog.Fatal("Invalid Format: no <REG> at start.")
			}
			f, err := processFiled(item.data, item.enums)
			if err != nil {
				return nil, err
			}
			if f.End >= jar.Width {
				clog.Fatal(fmt.Sprintf("Field offset %d is over bit-width %d.",
					f.End, jar.Width))
			}
			curReg.addFileds(f)
		}
	}

	if dbg.True {
		fmt.Println("----------------- after parse ---------------")
		fmt.Println(jar)
	}

	return jar, nil
}

func New(filename string) (*Jar, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// lines in file ---> tagItemSlice
	items, err := trim(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}

	// tagItemSlice -> Jar
	jar, err := parse(items)
	if err != nil {
		return nil, err
	}

	return jar, nil
}
