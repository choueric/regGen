package regjar

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/choueric/clog"
	"github.com/choueric/goutils"
)

type Field struct {
	Name      string
	Start     uint32
	End       uint32
	EnumVals  []uint32
	EnumNames []string
}

func (f *Field) String() string {
	var str bytes.Buffer
	fmt.Fprintf(&str, "%s: [%d:%d]", f.Name, f.Start, f.End)
	fmt.Fprintf(&str, " (")
	for i, v := range f.EnumVals {
		fmt.Fprintf(&str, "%d: %s, ", v, f.EnumNames[i])
	}
	fmt.Fprintf(&str, ")")
	return str.String()
}

// validate the field format like `name: offset (enumVal: enumName, enumVal:enumName)`
// which contains two parts:
//   - nameOffset, i.e. `name: offset`
//   - enums, i.e. `(enumVal: enumName, val:name)`. This part is optional
// If format is not correct, return (nil, false)
// Otherwise, return ({nameOffset str, enums str}, true)
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

func processFiledNameOffset(f *Field, nameOffset string) error {
	strs := strings.Split(nameOffset, ":")
	if len(strs) != 2 {
		return errors.New("Invalid Format: " + nameOffset)
	}

	f.Name = strings.ToUpper(strings.TrimSpace(strs[0]))
	offsetStr := strings.TrimSpace(strs[1])
	strs = strings.Split(offsetStr, "-")
	if len(strs) == 1 {
		offset, err := strconv.ParseInt(strs[0], 0, 32)
		if err != nil {
			clog.Error(nameOffset)
			return err
		}
		f.End = uint32(offset)
		f.Start = f.End
	} else if len(strs) == 2 {
		offset, err := strconv.ParseInt(strings.TrimSpace(strs[0]), 0, 32)
		if err != nil {
			clog.Error(nameOffset)
			return err
		}
		f.Start = uint32(offset)

		offset, err = strconv.ParseInt(strings.TrimSpace(strs[1]), 0, 32)
		if err != nil {
			clog.Error(nameOffset)
			return err
		}
		f.End = uint32(offset)

		if f.Start > f.End {
			f.Start, f.End = f.End, f.Start
		}
	}

	return nil
}

func processFiledEnums(f *Field, enumStr string) error {
	if enumStr == "" {
		return nil
	}

	nameStrs := []string{}
	valStrs := []string{}

	addPair := func(str string) {
		pair := strings.Split(str, ":")
		if len(pair) != 2 {
			clog.Fatal("Invalid field value format: " + enumStr)
		}
		valStrs = append(valStrs, strings.TrimSpace(pair[0]))
		nameStrs = append(nameStrs, strings.TrimSpace(pair[1]))
	}

	strs := strings.Split(enumStr, ",")
	if len(strs) == 1 { // may contain single pair
		addPair(enumStr)
	} else {
		for _, pairStr := range strs {
			addPair(pairStr)
		}
	}

	for i, ns := range nameStrs {
		if v, err := goutils.ParseUint(valStrs[i], 32); err != nil {
		} else {
			f.EnumVals = append(f.EnumVals, uint32(v))
			f.EnumNames = append(f.EnumNames, strings.ToUpper(ns))
		}
	}

	return nil
}

func processFiled(nameOffset, enums string) (*Field, error) {
	f := new(Field)

	if err := processFiledNameOffset(f, nameOffset); err != nil {
		return nil, err
	}

	if err := processFiledEnums(f, enums); err != nil {
		return nil, err
	}

	return f, nil
}
