package regjar

import (
	"bytes"
	"fmt"
)

type Reg struct {
	Name   string
	Offset uint64
	Fields []*Field
}

func (r *Reg) String() string {
	var str bytes.Buffer
	fmt.Fprintf(&str, "\"%s\", %#x\n", r.Name, r.Offset)
	for _, f := range r.Fields {
		fmt.Fprintln(&str, "   ", f)
	}
	return str.String()
}

func (r *Reg) addFileds(f ...*Field) {
	r.Fields = append(r.Fields, f...)
}

func processReg(line string) (*Reg, error) {
	name, offset, err := parseTagNameOffset(line)
	if err != nil {
		return nil, err
	}

	return &Reg{name, offset, make([]*Field, 0)}, nil
}
