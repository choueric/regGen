package regjar

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/choueric/clog"
	"github.com/choueric/regGen/dbg"
)

const (
	tag_chip = iota
	tag_width
	tag_module
	tag_reg
	tag_field
	tag_comment
	tag_other
)

type tagItem struct {
	tag   int
	data  string
	enums string
}

func (item *tagItem) String() string {
	switch item.tag {
	case tag_chip:
		return fmt.Sprintf("[  CHIP ] %s", item.data)
	case tag_width:
		return fmt.Sprintf("[ WIDTH ] %s", item.data)
	case tag_module:
		return fmt.Sprintf("[  MOD  ] %s", item.data)
	case tag_reg:
		return fmt.Sprintf("[  REG  ] %s", item.data)
	case tag_comment:
		return fmt.Sprintf("[ COMNT ] %s", item.data)
	case tag_field:
		if item.enums == "" {
			return fmt.Sprintf("[ FIELD ] %s", item.data)
		} else {
			return fmt.Sprintf("[ FIELD ] %s (%s)", item.data, item.enums)
		}
	case tag_other:
		return fmt.Sprintf("[ OTHER ] %s", item.data)
	default:
		clog.Fatal("Unkonw type: " + item.data)
		return ""
	}
}

type tagItemSlice []*tagItem

func (s tagItemSlice) String() string {
	var str bytes.Buffer
	for _, i := range s {
		fmt.Fprintln(&str, i)
	}
	return str.String()
}

func (s *tagItemSlice) addTagItems(v ...*tagItem) {
	p := *s
	p = append(p, v...)
	*s = p
}

func newTagItemSlice() tagItemSlice {
	return tagItemSlice(make([]*tagItem, 0))
}

func parseTagItem(line string) *tagItem {
	sLine := strings.TrimSpace(line)

	// comment
	if m, _ := regexp.MatchString(`\s*#`, sLine); m {
		return &tagItem{tag: tag_comment, data: sLine}
	}

	// other
	strs := strings.Split(sLine, ":")
	if len(strs) == 1 {
		return &tagItem{tag: tag_other, data: sLine}
	}

	// valid tag line
	tagStr := strings.ToUpper(strs[0])
	if strings.Contains(tagStr, "<CHIP>") {
		return &tagItem{tag: tag_chip, data: sLine}
	} else if strings.Contains(tagStr, "<WIDTH>") {
		return &tagItem{tag: tag_width, data: sLine}
	} else if strings.Contains(tagStr, "<MODULE>") {
		return &tagItem{tag: tag_module, data: sLine}
	} else if strings.Contains(tagStr, "<REG>") {
		return &tagItem{tag: tag_reg, data: sLine}
	} else {
		if strs, ok := validField(sLine); ok {
			return &tagItem{tag: tag_field, data: strs[0], enums: strs[1]}
		} else {
			if len(line) != 0 {
				clog.Fatal("Invalid Format: [" + line + "]")
			}
		}
	}
	return nil
}

func newTagItem(line string) *tagItem {
	item := parseTagItem(line)
	if dbg.True {
		fmt.Println(item)
	}
	return item
}
