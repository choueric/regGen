package format

import (
	"errors"
	"io"

	"github.com/choueric/regGen/regjar"
)

type Formatter interface {
	FormatLicense(w io.Writer, license string)
	FormatRegJar(w io.Writer, jar *regjar.Jar)
}

func New(fmt string) (Formatter, error) {
	switch fmt {
	case "cmacro":
		return new(cmacroFormat), nil
	default:
		return nil, errors.New("format: invalid output format: " + fmt)
	}
}
