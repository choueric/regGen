package regjar

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/choueric/clog"
)

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
