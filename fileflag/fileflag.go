package fileflag

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"

	"github.com/choueric/goutils"
)

type FileFlag struct {
	set *flag.FlagSet
}

func New(filepath string) *FileFlag {
	return &FileFlag{set: flag.NewFlagSet(filepath, flag.ExitOnError)}
}

func (ff *FileFlag) FlagSet() *flag.FlagSet {
	return ff.set
}

// Parse first parses the commandline if it is supplied, otherwise parse
// the flag file whose path is filepath from New()
func (ff *FileFlag) Parse() error {
	if len(os.Args) > 1 {
		ff.set.Parse(os.Args[1:])
		return nil
	}

	if isExist, err := goutils.IsFileExist(ff.set.Name()); err != nil {
		return err
	} else {
		if !isExist {
			ff.set.Parse(os.Args[1:])
			return nil
		}
	}

	f, err := os.Open(ff.set.Name())
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	args := make([]string, 0)
	for {
		line, err := goutils.ReadLine(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		args = append(args, strings.Fields(line)...)
	}

	ff.set.Parse(args)
	return nil
}
