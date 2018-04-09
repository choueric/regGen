package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/choueric/clog"
)

const (
	LINE_TAG_CHIP = iota
	LINE_TAG_REG
	LINE_TAG_FIELD
	LINE_TAG_COMMENT
	LINE_TAG_OTHER
)

type item struct {
	data string
	tag  int
}

func printTrimItems(w io.Writer, items []item) {
	for _, i := range items {
		switch i.tag {
		case LINE_TAG_CHIP:
			fmt.Fprintln(w, "[  CHIP ]", i.data)
		case LINE_TAG_REG:
			fmt.Fprintln(w, "[  REG  ]", i.data)
		case LINE_TAG_FIELD:
			fmt.Fprintln(w, "[ FIELD ]", i.data)
		}
	}
}

type field struct {
	name  string
	start uint32
	end   uint32
}

func (f *field) String() string {
	return fmt.Sprintf("%s: [%d:%d]", f.name, f.start, f.end)
}

type reg struct {
	name   string
	offset uint64
	fields []*field
}

func (r *reg) String() string {
	return fmt.Sprintf("\"%s\", %#x", r.name, r.offset)
}

type regMap struct {
	chip string
	regs []*reg
}

func (rm *regMap) String() string {
	var str bytes.Buffer
	fmt.Fprintf(&str, "CHIP: \"%s\"\n", rm.chip)
	for _, r := range rm.regs {
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

func get_line_tag(line string) int {
	if strings.Contains(line, "<CHIP>") || strings.Contains(line, "<chip>") {
		if debug {
			fmt.Println("[  CHIP ]", line)
		}
		return LINE_TAG_CHIP
	} else if strings.Contains(line, "<REG>") || strings.Contains(line, "<reg>") {
		if debug {
			fmt.Println("[  REG  ]", line)
		}
		return LINE_TAG_REG
	} else if strings.Contains(line, ":") {
		if debug {
			fmt.Println("[ FIELD ]", line)
		}
		return LINE_TAG_FIELD
	} else if strings.Contains(line, "#") {
		if debug {
			fmt.Println("[ COMET ]", line)
		}
		return LINE_TAG_COMMENT
	} else {
		if debug {
			fmt.Println("[ OTHER ]", line)
		}
		return LINE_TAG_OTHER
	}
}

func processChip(line string) (string, error) {
	strs := strings.Split(line, ":")
	if len(strs) != 2 {
		return "", errors.New("Invalid Format: " + line)
	}
	return strings.TrimSpace(strs[1]), nil
}

func processReg(line string) (*reg, error) {
	r := &reg{}
	strs := strings.Split(line, ":")
	if len(strs) != 2 {
		return nil, errors.New("Invalid Format: " + line)
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

func processFiled(line string) (*field, error) {
	f := &field{}
	strs := strings.Split(line, ":")
	if len(strs) != 2 {
		return nil, errors.New("Invalid Format: " + line)
	}

	f.name = strings.TrimSpace(strs[0])
	offsetStr := strings.TrimSpace(strs[1])
	strs = strings.Split(offsetStr, "-")
	if len(strs) == 1 {
		offset, err := strconv.ParseInt(strs[0], 0, 32)
		if err != nil {
			clog.Error(line)
			return nil, err
		}
		f.end = uint32(offset)
		f.start = f.end
	} else if len(strs) == 2 {
		offset, err := strconv.ParseInt(strings.TrimSpace(strs[0]), 0, 32)
		if err != nil {
			clog.Error(line)
			return nil, err
		}
		f.start = uint32(offset)

		offset, err = strconv.ParseInt(strings.TrimSpace(strs[1]), 0, 32)
		if err != nil {
			clog.Error(line)
			return nil, err
		}
		f.end = uint32(offset)

		if f.start > f.end {
			f.start, f.end = f.end, f.start
		}

	}

	return f, nil
}

func trim(reader *bufio.Reader) ([]item, error) {
	items := make([]item, 0)
	for {
		line, err := readLine(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		t := get_line_tag(line)
		switch t {
		case LINE_TAG_CHIP:
			fallthrough
		case LINE_TAG_REG:
			fallthrough
		case LINE_TAG_FIELD:
			items = append(items, item{data: strings.TrimSpace(line), tag: t})
		case LINE_TAG_COMMENT:
		default:
			if len(line) != 0 {
				clog.Warn("Invalid Format:", line)
			}

		}
	}

	if debug {
		fmt.Println("----------------- after trim ---------------")
		printTrimItems(os.Stdout, items)
	}

	return items, nil
}

func parse(rm *regMap, items []item) error {
	var curReg *reg
	for _, item := range items {
		switch item.tag {
		case LINE_TAG_CHIP:
			chip, err := processChip(item.data)
			if err != nil {
				return err
			}
			rm.chip = chip
		case LINE_TAG_REG:
			r, err := processReg(item.data)
			if err != nil {
				return err
			}
			rm.regs = append(rm.regs, r)
			curReg = r
		case LINE_TAG_FIELD:
			if curReg == nil {
				return errors.New("Invalid Format: no <REG> at start")
			}
			f, err := processFiled(item.data)
			if err != nil {
				return err
			}
			curReg.fields = append(curReg.fields, f)
		}
	}

	if debug {
		fmt.Println("----------------- after parse ---------------")
		fmt.Println(rm)
	}

	return nil
}

func loadRegs(rm *regMap, reader *bufio.Reader) error {
	items, err := trim(reader)
	if err != nil {
		return err
	}

	err = parse(rm, items)
	if err != nil {
		return err
	}

	return nil
}

func (rm *regMap) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = loadRegs(rm, bufio.NewReader(f))
	if err != nil {
		return err
	}

	return nil
}

func (rm *regMap) Output(w io.Writer, format string) error {
	f, ok := outputFormat[format]
	if !ok {
		return errors.New("Invalid format: " + format)
	}

	if debug {
		fmt.Println("----------------- format output ---------------")
	}
	f(rm, w)
	return nil
}
