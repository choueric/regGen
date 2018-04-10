package main

import (
	"errors"
	"fmt"
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
