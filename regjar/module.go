package regjar

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/choueric/clog"
)

type Module struct {
	Name string
	Regs []*Reg
}

func defaultModule() *Module {
	return &Module{"default", make([]*Reg, 0)}
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
	strs := strings.Split(line, ":")
	if len(strs) != 2 {
		clog.Fatal("Invalid Format: [" + line + "]")
	}

	return &Module{strings.TrimSpace(strs[1]), make([]*Reg, 0)}, nil
}
