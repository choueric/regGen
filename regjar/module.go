package regjar

import (
	"bytes"
	"fmt"
)

type Module struct {
	Name     string
	BaseAddr uint64
	Regs     []*Reg
}

func defaultModule() *Module {
	return &Module{"default", 0, make([]*Reg, 0)}
}

func (mod *Module) String() string {
	var str bytes.Buffer
	fmt.Fprintf(&str, "MODULE: \"%s\"\n", mod.Name)
	for _, r := range mod.Regs {
		fmt.Fprint(&str, r)
	}
	return str.String()
}

func (mod *Module) addRegs(v ...*Reg) {
	mod.Regs = append(mod.Regs, v...)
}

func processModule(line string) (*Module, error) {
	name, addr, err := parseTagNameOffset(line)
	if err != nil {
		return nil, err
	}

	return &Module{name, addr, make([]*Reg, 0)}, nil
}
