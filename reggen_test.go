package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

const (
	inputStr = `# comment
	<chip>:si5324
	<REG>[Control]: 0
	BYPASS_REG: 1
	CKOUT_ALWAYS_ON: 5
	FREE_RUN: 7 - 6 

<REG>: 1
  ck_prior1 : 0-1 
  ck_prior2; 2-3

<REG>: 0x10
BWSEL_REG: 4-7`
	trimedStr = `[  CHIP ] <chip>:si5324
[  REG  ] <REG>[Control]: 0
[ FIELD ] BYPASS_REG: 1
[ FIELD ] CKOUT_ALWAYS_ON: 5
[ FIELD ] FREE_RUN: 7 - 6
[  REG  ] <REG>: 1
[ FIELD ] ck_prior1 : 0-1
[  REG  ] <REG>: 0x10
`
	parsedStr = `CHIP: "si5324"
"Control", 0x0
    BYPASS_REG: [1:1]
    CKOUT_ALWAYS_ON: [5:5]
    FREE_RUN: [6:7]
"1", 0x1
    ck_prior1: [0:1]
"16", 0x10
`
	formatCStr = `// Registers of si5324
#define BIT(x) (1 << (x))
#define MASK(a, b) (((uint8_t)-1 >> (7-(b))) & ~((1U<<(a))-1))

#define REG_CONTROL 0x0
	#define REG_BYPASS_REG_BIT BIT(1)
	#define REG_CKOUT_ALWAYS_ON_BIT BIT(5)
	#define REG_FREE_RUN_MSK MASK(6, 7)

#define REG_1 0x1
	#define REG_CK_PRIOR1_MSK MASK(0, 1)

#define REG_16 0x10
`
)

func print_string_mismatch(a, b []byte) {
	if len(a) != len(b) {
		fmt.Println("length doesn't match:", len(a), len(b))
		return
	}
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			fmt.Println(i, string(a[i]), string(b[i]))
			break
		}
	}
}

func Test_trim(t *testing.T) {
	items, err := trim(bufio.NewReader(strings.NewReader(inputStr)))
	if err != nil {
		t.Error(err)
	}

	var result bytes.Buffer
	printTrimItems(&result, items)
	if strings.Compare(result.String(), trimedStr) != 0 {
		print_string_mismatch(result.Bytes(), []byte(trimedStr))
		t.Error("trim fail!")
	}
}

func Test_parse(t *testing.T) {
	items, err := trim(bufio.NewReader(strings.NewReader(inputStr)))
	if err != nil {
		t.Error(err)
	}

	rm := regMap{}
	err = parse(&rm, items)
	if err != nil {
		t.Error(err)
	}

	var result bytes.Buffer
	fmt.Fprint(&result, &rm)
	if strings.Compare(result.String(), parsedStr) != 0 {
		print_string_mismatch(result.Bytes(), []byte(parsedStr))
		t.Error("parse fail!")
	}
}

func Test_formatToC(t *testing.T) {
	items, err := trim(bufio.NewReader(strings.NewReader(inputStr)))
	if err != nil {
		t.Error(err)
	}

	rm := regMap{}
	err = parse(&rm, items)
	if err != nil {
		t.Error(err)
	}

	var result bytes.Buffer
	formatToC(&rm, &result)
	if strings.Compare(result.String(), formatCStr) != 0 {
		print_string_mismatch(result.Bytes(), []byte(formatCStr))
		t.Error("parse fail!")
	}
}
