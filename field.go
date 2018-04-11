package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/choueric/clog"
)

type field struct {
	name    string
	start   uint32
	end     uint32
	valData []uint32
	valName []string
}

func (f *field) String() string {
	var str bytes.Buffer
	fmt.Fprintf(&str, "%s: [%d:%d]", f.name, f.start, f.end)
	fmt.Fprintf(&str, " (")
	for i, v := range f.valData {
		fmt.Fprintf(&str, "%d: %s, ", v, f.valName[i])
	}
	fmt.Fprintf(&str, ")")
	return str.String()
}

// `0b000111` to uint32
func parseBinStr(str string) (uint32, error) {
	if v, err := strconv.ParseUint(str[2:], 2, 32); err != nil {
		return 0, err
	} else {
		return uint32(v), nil
	}
}

// base 16`0x`, base 8`0` and base 10 to uint32
func parseOtherBaseStr(str string) (uint32, error) {
	if v, err := strconv.ParseUint(str, 0, 32); err != nil {
		return 0, err
	} else {
		return uint32(v), nil
	}
}

// validate the field format like `name: offset (val: vname, val:vname)`
// which contains two parts:
//   - nameOffset, i.e. `name: offset`
//   - values, i.e. `(val: vname, val:vname)`. This part is optional
// If format is not correct, return (nil, false)
// Otherwise, return ({str of nameOffset, str of values}, true)
func validField(line string) ([]string, bool) {
	if !strings.Contains(line, ":") {
		return nil, false
	}

	if strings.Contains(line, "(") || strings.Contains(line, ")") {
		m := regexp.MustCompile(`.*\((.+)\)`)
		strs := m.FindStringSubmatch(line)
		if len(strs) == 2 {
			subs := strings.Split(line, "(")
			return []string{strings.TrimSpace(subs[0]), strs[1]}, true
		} else {
			return nil, false
		}
	}

	return []string{line, ""}, true
}

func processFiledNameOffset(f *field, nameOffset string) error {
	strs := strings.Split(nameOffset, ":")
	if len(strs) != 2 {
		return errors.New("Invalid Format: " + nameOffset)
	}

	f.name = strings.TrimSpace(strs[0])
	offsetStr := strings.TrimSpace(strs[1])
	strs = strings.Split(offsetStr, "-")
	if len(strs) == 1 {
		offset, err := strconv.ParseInt(strs[0], 0, 32)
		if err != nil {
			clog.Error(nameOffset)
			return err
		}
		f.end = uint32(offset)
		f.start = f.end
	} else if len(strs) == 2 {
		offset, err := strconv.ParseInt(strings.TrimSpace(strs[0]), 0, 32)
		if err != nil {
			clog.Error(nameOffset)
			return err
		}
		f.start = uint32(offset)

		offset, err = strconv.ParseInt(strings.TrimSpace(strs[1]), 0, 32)
		if err != nil {
			clog.Error(nameOffset)
			return err
		}
		f.end = uint32(offset)

		if f.start > f.end {
			f.start, f.end = f.end, f.start
		}
	}

	return nil
}

func processFiledValues(f *field, valStr string) error {
	if valStr == "" {
		return nil
	}

	nameStrs := []string{}
	valStrs := []string{}

	addPair := func(str string) {
		pair := strings.Split(str, ":")
		if len(pair) != 2 {
			clog.Fatal("Invalid field value format: " + valStr)
		}
		valStrs = append(valStrs, strings.TrimSpace(pair[0]))
		nameStrs = append(nameStrs, strings.TrimSpace(pair[1]))
	}

	strs := strings.Split(valStr, ",")
	if len(strs) == 1 { // may contain single pair
		addPair(valStr)
	} else {
		for _, pairStr := range strs {
			addPair(pairStr)
		}
	}

	parseFunc := map[bool]func(string) (uint32, error){
		true:  parseBinStr,
		false: parseOtherBaseStr,
	}

	for i, ns := range nameStrs {
		matched, err := regexp.MatchString(`^0b[01]*`, valStrs[i])
		if err != nil {
			clog.Fatal(err)
		}

		if v, err := parseFunc[matched](valStrs[i]); err != nil {
			clog.Fatal(err)
		} else {
			f.valData = append(f.valData, v)
			f.valName = append(f.valName, ns)
		}
	}

	return nil
}

func processFiled(nameOffset, valStr string) (*field, error) {
	f := &field{}

	if err := processFiledNameOffset(f, nameOffset); err != nil {
		return nil, err
	}

	if err := processFiledValues(f, valStr); err != nil {
		return nil, err
	}

	return f, nil
}
