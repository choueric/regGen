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
	tag         int
	data        string
	fieldValStr string // only for filed line
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
		return fmt.Sprintf("[ FIELD ] %s", item.data)
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
			fmt.Fprintf(&str, "[ FIELD ] %s (%s)\n", i.data, i.fieldValStr)
		}
	}
	return str.String()
}

type reg struct {
	name   string
	offset uint64
	fields []*field
}

func (r *reg) String() string {
	return fmt.Sprintf("\"%s\", %#x", r.name, r.offset)
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

func readLine(r *bufio.Reader) (string, error) {
	str, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	str = strings.Trim(str, "\r\n")
	return str, nil
}

func tagItemNew(line string) (item *tagItem) {
	sLine := strings.TrimSpace(line)
	if strings.Contains(sLine, "<CHIP>") || strings.Contains(sLine, "<chip>") {
		item = &tagItem{tag: TAG_CHIP, data: sLine}
	} else if strings.Contains(sLine, "<REG>") || strings.Contains(sLine, "<reg>") {
		item = &tagItem{tag: TAG_REG, data: sLine}
	} else if m, _ := regexp.MatchString(`\s*#`, sLine); m {
		item = &tagItem{tag: TAG_COMMENT, data: sLine}
	} else {
		if strs, ok := validField(sLine); ok {
			item = &tagItem{tag: TAG_FIELD, data: strs[0], fieldValStr: strs[1]}
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
	r := &reg{}
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
	items := tagItemSlice(make([]*tagItem, 0))
	for {
		line, err := readLine(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		item := tagItemNew(line)
		if err != nil {
			clog.Fatal(err)
		} else {
			if item != nil && item.tag != TAG_COMMENT {
				items = append(items, item)
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
	jar := &regJar{}
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
			jar.regs = append(jar.regs, r)
			curReg = r
		case TAG_FIELD:
			if curReg == nil {
				clog.Fatal("Invalid Format: no <REG> at start")
			}
			f, err := processFiled(item.data, item.fieldValStr)
			if err != nil {
				return nil, err
			}
			curReg.fields = append(curReg.fields, f)
		}
	}

	if debug {
		fmt.Println("----------------- after parse ---------------")
		fmt.Println(jar)
	}

	return jar, nil
}

func regJarNew(filename string) (*regJar, error) {
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
