# regGen

It is a tool written in Golang to automatically generate code of registers and
fields in registers of a chip by reading a register-description file.

Now it firstly supports the auto-generation of C format because it's the most
application scence.

# usage

To see the usage information:

```sh
$ reggen -h
Usage of reggen:
  -d    enable debug
  -f string
        output format type. [c] (default "c")
  -i string
        input file. (default "input.regs")
```

So to output what you want, just:

```sh
$ reggen -i input.regs > regs.h
```

By redirection, you will get `regs.h` containing the macros about registers and
fields.

# example

The content of `input.regs` is:

```
# si5324 is a clock generator
<chip>:si5324

<REG>[Control]: 0
        BYPASS_REG: 1
        CKOUT_ALWAYS_ON: 5
        FREE_RUN: 7 - 6 

<REG>: 1
  ck_prior1 : 0-1 
  ck_prior2; 2-3

<REG>: 0x10
BWSEL_REG: 4-7
```

And the output of C format will be:

```c
#pragma once

#ifndef BIT
#define BIT(x) (1 << (x))
#endif

// ONLY for _8bit-width_ register
#define MASK(a, b) (((uint8_t)-1 >> (7-(b))) & ~((1U<<(a))-1))

// Registers of si5324

#define REG_CONTROL 0x0 // 0
	#define REG_BYPASS_REG_BIT BIT(1)
	#define REG_CKOUT_ALWAYS_ON_BIT BIT(5)
	#define REG_FREE_RUN_MSK MASK(6, 7)
	#define REG_FREE_RUN_VAL(rv) (((rv) & REG_FREE_RUN_MSK) >> 6)
	#define REG_FREE_RUN_SFT(v) (((v) & REG_FREE_RUN_MSK) << 6)

#define REG_1 0x1 // 1
	#define REG_CK_PRIOR1_MSK MASK(0, 1)
	#define REG_CK_PRIOR1_VAL(rv) ((rv) & REG_CK_PRIOR1_MSK)
	#define REG_CK_PRIOR1_SFT(v) ((v) & REG_CK_PRIOR1_MSK)

#define REG_16 0x10 // 16
	#define REG_BWSEL_REG_MSK MASK(4, 7)
	#define REG_BWSEL_REG_VAL(rv) (((rv) & REG_BWSEL_REG_MSK) >> 4)
	#define REG_BWSEL_REG_SFT(v) (((v) & REG_BWSEL_REG_MSK) << 4)

```

# input format

The above example shows the input file's basic format which contains only a few
keywords and follows simple and loose rules. All keywords are case-insensitive.

- `COMMENT`: A line is treated as comment line if starting with `#`.
  
  Example: `# si5324 is a clock generator`.
- `CHIP`: Specify the chip's name. It is optional and name is empty.
  Example: `<chip>:si5324`, where `si5324` is the name.
  
- `REG`: : It defines the name (optional) and offset of a register and generally
  followed with serveral `FIELD` lines which defines this register's fields.
  
  Example: `<REG>[Control]: 0`, where the register's name is `Control` and offset is 0.
  
- `FIELD`: As says above, `FIELD` lines contains information of a register's
  filed, including name, offset and values (optional) in the format of
  `name: offset (value: valueName,)` that contains two parts. The `name-offset`
  part is mandatory contrasting to the `values` part. Below explains more details
  about `FIELD` line.
  
  Example: `ck_prior1 : 0-1`, where the field's name is `ck_prior1`, offset
  range is `0-1` and there is no `values` part.
  
## `name-offset` part

This part of `FIELD` line contains contains field's name and offset. The offset
are either bit or range which use `-` to connect. The start bit and end bit can
in either sides of `-`.
  
## `values` part

This part of `FIELD` line is optional and gives the value range this field can
have, which is usually used when the field's offset is not only one bit.

Embraced by `(` and `)`, the values are paired with `:`, where the value number
is at left and value name right. Value pairs are separated with `,`.

For example, the chip spec have one filed like:

```
4:3 | VALTIME[2]
00: 2ms
01: 100ms
10: 200ms
11: 13 seconds
```

Then the `FIELD` line should be one of lines below:
```
VALTIME: 4-3 (0b00: 2ms, 0b01: 100ms, 0b10: 200ms, 0b11: 13s)
VALTIME: 4-3 (0: 2ms, 1: 100ms, 2: 200ms, 3: 13s)
```

As the example shows, the value number can be decimal and binary. Actually, it
support four formats:

- dec, like 10 = 10
- bin, like 0b11 = 3
- oct, like 012 = 10
- hex, like 0x1f = 31

## spec

A formal-like sepcification is:

```
<chip>: chipName

<reg>[regName]: offset
filed_1: bit
filed_2: startBit - endBit
filed_3: endBit - startBit (val1: val1Name, val2: val2Name)
```

# output format

Right now, the only one output format is the C format which uses `#define`
macros to represent all registers and fields and capitalize all their names.

It handle two types of field: `bit` and `mask`:
The output formats for each type of `FIELD` line are different:
- for `bit` type, it only generates one single macro `REG_XXX_BIT`
- for `range` type, it outputs macros below:
	- `REG_XXX_MSK`: the mask for this field. For example, `REG_XXX_MSK = 0x18`
	  for `3-4`.
	- `REG_XXX_VAL(rv)`: this macro can get this field's value from register
	  value `rv`. For example, `REG_XXX_VAL(0x10) = 2` for `3-4`.
	- `REG_XXX_SFT(v)`: it shifts the field's value `v` to the correct offset.
	  For example, `REG_XXX_SFT(1) = 0x08` for `3-4`.
	- `REG_XXX_VALNAME VALNUM`: if the line contains `values` part. In the case
	  of the VALTIME example above, the output is:
	  ```
	  #define REG_VALTIME_2MS 	0 	// 0b0	0x0
	  #define REG_VALTIME_100MS 	1 	// 0b1	0x1
	  #define REG_VALTIME_200MS 	2 	// 0b10	0x2
	  #define REG_VALTIME_13S 	3 	// 0b11	0x3
	  ```
