package regjar

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
	"github.com/choueric/regGen/dbg"
)

const (
	tag_chip = iota
	tag_reg
	tag_field
	tag_comment
	tag_other
)

type tagItem struct {
	tag   int
	data  string
	enums string
}

func (item *tagItem) String() string {
	switch item.tag {
	case tag_chip:
		return fmt.Sprintf("[  CHIP ] %s", item.data)
	case tag_reg:
		return fmt.Sprintf("[  REG  ] %s", item.data)
	case tag_comment:
		return fmt.Sprintf("[ COMNT ] %s", item.data)
	case tag_field:
		if item.enums == "" {
			return fmt.Sprintf("[ FIELD ] %s", item.data)
		} else {
			return fmt.Sprintf("[ FIELD ] %s (%s)", item.data, item.enums)
		}
	case tag_other:
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
		case tag_chip:
			fmt.Fprintln(&str, "[  CHIP ]", i.data)
		case tag_reg:
			fmt.Fprintln(&str, "[  REG  ]", i.data)
		case tag_field:
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

type Reg struct {
	Name   string
	Offset uint64
	Fields []*Field
}

func (r *Reg) String() string {
	return fmt.Sprintf("\"%s\", %#x", r.Name, r.Offset)
}

func (r *Reg) addFileds(f ...*Field) {
	r.Fields = append(r.Fields, f...)
}

type Jar struct {
	Chip string
	Regs []*Reg
}

func (jar *Jar) String() string {
	var str bytes.Buffer
	fmt.Fprintf(&str, "CHIP: \"%s\"\n", jar.Chip)
	for _, r := range jar.Regs {
		fmt.Fprintln(&str, r)
		for _, f := range r.Fields {
			fmt.Fprintln(&str, "   ", f)
		}
	}
	return str.String()
}

func (jar *Jar) addRegs(v ...*Reg) {
	jar.Regs = append(jar.Regs, v...)
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
		item = &tagItem{tag: tag_chip, data: sLine}
	} else if strings.Contains(sLine, "<REG>") || strings.Contains(sLine, "<reg>") {
		item = &tagItem{tag: tag_reg, data: sLine}
	} else if m, _ := regexp.MatchString(`\s*#`, sLine); m {
		item = &tagItem{tag: tag_comment, data: sLine}
	} else {
		if strs, ok := validField(sLine); ok {
			item = &tagItem{tag: tag_field, data: strs[0], enums: strs[1]}
		} else {
			if len(line) != 0 {
				clog.Fatal("Invalid Format: [" + line + "]")
			}
		}
	}

	if dbg.True {
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

func processReg(line string) (*Reg, error) {
	r := new(Reg)
	strs := strings.Split(line, ":")
	if len(strs) != 2 {
		clog.Fatal("Invalid Format: [" + line + "]")
	} else {
		offset, err := strconv.ParseInt(strings.TrimSpace(strs[1]), 0, 64)
		if err != nil {
			clog.Error(line)
			return nil, err
		}
		r.Offset = uint64(offset)
	}

	a := strings.IndexByte(line, '[')
	b := strings.IndexByte(line, ']')
	if a != -1 && b != -1 {
		r.Name = strings.TrimSpace(line[a+1 : b])
	} else {
		r.Name = strconv.FormatUint(r.Offset, 10)
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
			if item != nil && item.tag != tag_comment {
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
	jar := new(Jar)
	for _, item := range items {
		switch item.tag {
		case tag_chip:
			chip, err := processChip(item.data)
			if err != nil {
				return nil, err
			}
			jar.Chip = chip
		case tag_reg:
			r, err := processReg(item.data)
			if err != nil {
				return nil, err
			}
			jar.addRegs(r)
			curReg = r
		case tag_field:
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
