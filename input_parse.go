package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/choueric/clog"
)

const (
	TAG_CHIP = iota
	TAG_REG
	TAG_FIELD
	TAG_COMMENT
	TAG_OTHER
)

type tagItem struct {
	tag   int
	data  string
	enums string
}

func (item *tagItem) String() string {
	switch item.tag {
	case TAG_CHIP:
		return fmt.Sprintf("[  CHIP ] %s", item.data)
	case TAG_REG:
		return fmt.Sprintf("[  REG  ] %s", item.data)
	case TAG_COMMENT:
		return fmt.Sprintf("[ COMNT ] %s", item.data)
	case TAG_FIELD:
		if item.enums == "" {
			return fmt.Sprintf("[ FIELD ] %s", item.data)
		} else {
			return fmt.Sprintf("[ FIELD ] %s (%s)", item.data, item.enums)
		}
	case TAG_OTHER:
		return fmt.Sprintf("[ OTHER ] %s", item.data)
	default:
		clog.Fatal("Unkonw type: " + item.data)
		return ""
	}
}

type tagItemSlice []*tagItem

func (s tagItemSlice) String() string {
	var str bytes.Buffer
	for _, i := range s {
		switch i.tag {
		case TAG_CHIP:
			fmt.Fprintln(&str, "[  CHIP ]", i.data)
		case TAG_REG:
			fmt.Fprintln(&str, "[  REG  ]", i.data)
		case TAG_FIELD:
			fmt.Fprintf(&str, "[ FIELD ] %s (%s)\n", i.data, i.enums)
		}
	}
	return str.String()
}

func (s *tagItemSlice) addTagItems(v ...*tagItem) {
	p := *s
	p = append(p, v...)
	*s = p
}

func newTagItemSlice() tagItemSlice {
	return tagItemSlice(make([]*tagItem, 0))
}

type reg struct {
	name   string
	offset uint64
	fields []*field
}

func (r *reg) String() string {
	return fmt.Sprintf("\"%s\", %#x", r.name, r.offset)
}

func (r *reg) addFileds(f ...*field) {
	r.fields = append(r.fields, f...)
}

type regJar struct {
	chip string
	regs []*reg
}

func (jar *regJar) String() string {
	var str bytes.Buffer
	fmt.Fprintf(&str, "CHIP: \"%s\"\n", jar.chip)
	for _, r := range jar.regs {
		fmt.Fprintln(&str, r)
		for _, f := range r.fields {
			fmt.Fprintln(&str, "   ", f)
		}
	}
	return str.String()
}

func (jar *regJar) addRegs(v ...*reg) {
	jar.regs = append(jar.regs, v...)
}

func readLine(r *bufio.Reader) (string, error) {
	str, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	str = strings.Trim(str, "\r\n")
	return str, nil
}

func newTagItem(line string) (item *tagItem) {
	sLine := strings.TrimSpace(line)
	if strings.Contains(sLine, "<CHIP>") || strings.Contains(sLine, "<chip>") {
		item = &tagItem{tag: TAG_CHIP, data: sLine}
	} else if strings.Contains(sLine, "<REG>") || strings.Contains(sLine, "<reg>") {
		item = &tagItem{tag: TAG_REG, data: sLine}
	} else if m, _ := regexp.MatchString(`\s*#`, sLine); m {
		item = &tagItem{tag: TAG_COMMENT, data: sLine}
	} else {
		if strs, ok := validField(sLine); ok {
			item = &tagItem{tag: TAG_FIELD, data: strs[0], enums: strs[1]}
		} else {
			if len(line) != 0 {
				clog.Fatal("Invalid Format: [" + line + "]")
			}
		}
	}

	if debug {
		fmt.Println(item)
	}

	return
}

func processChip(line string) (string, error) {
	strs := strings.Split(line, ":")
	if len(strs) != 2 {
		clog.Fatal("Invalid Format: [" + line + "]")
	}
	return strings.TrimSpace(strs[1]), nil
}

func processReg(line string) (*reg, error) {
	r := new(reg)
	strs := strings.Split(line, ":")
	if len(strs) != 2 {
		clog.Fatal("Invalid Format: [" + line + "]")
	} else {
		offset, err := strconv.ParseInt(strings.TrimSpace(strs[1]), 0, 64)
		if err != nil {
			clog.Error(line)
			return nil, err
		}
		r.offset = uint64(offset)
	}

	a := strings.IndexByte(line, '[')
	b := strings.IndexByte(line, ']')
	if a != -1 && b != -1 {
		r.name = strings.TrimSpace(line[a+1 : b])
	} else {
		r.name = strconv.FormatUint(r.offset, 10)
	}
	return r, nil
}

func trim(reader *bufio.Reader) (tagItemSlice, error) {
	items := newTagItemSlice()
	for {
		line, err := readLine(reader)
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
			if item != nil && item.tag != TAG_COMMENT {
				items.addTagItems(item)
			}
		}
	}

	if debug {
		fmt.Println("----------------- after trim ---------------")
		fmt.Println(items)
	}

	return items, nil
}

func parse(items tagItemSlice) (*regJar, error) {
	var curReg *reg
	jar := new(regJar)
	for _, item := range items {
		switch item.tag {
		case TAG_CHIP:
			chip, err := processChip(item.data)
			if err != nil {
				return nil, err
			}
			jar.chip = chip
		case TAG_REG:
			r, err := processReg(item.data)
			if err != nil {
				return nil, err
			}
			jar.addRegs(r)
			curReg = r
		case TAG_FIELD:
			if curReg == nil {
				clog.Fatal("Invalid Format: no <REG> at start")
			}
			f, err := processFiled(item.data, item.enums)
			if err != nil {
				return nil, err
			}
			curReg.addFileds(f)
		}
	}

	if debug {
		fmt.Println("----------------- after parse ---------------")
		fmt.Println(jar)
	}

	return jar, nil
}

func newRegJar(filename string) (*regJar, error) {
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

	// tagItemSlice -> regJar
	jar, err := parse(items)
	if err != nil {
		return nil, err
	}

	return jar, nil
}
