package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/choueric/clog"
)

type field struct {
	name  string
	start uint32
	end   uint32
}

func (f *field) String() string {
	return fmt.Sprintf("%s: [%d:%d]", f.name, f.start, f.end)
}

// validate the field format like `name: offset (val: v, val:v)`
// which contains two parts:
//   - nameOffset, i.e. `name: offset`
//   - values, i.e. `(val: v, val:v)`. This part is optional
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
			if debug {
				for _, s := range strs {
					fmt.Println("  ", s)
				}
			}
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
