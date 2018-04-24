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
# This is a simple regs file
<chip>:simpleChip

<REG>[Control]: 0
        BYPASS: 1
        FREE_RUN: 7 - 6 

<REG>: 0x2
type: 4-5 (0b00: client, 0x01: server, 2: route, 3: peer)
```

And the output of C format will be:

```c
#pragma once

#ifndef BIT
#define BIT(x) (1 << (x))
#endif

// ONLY for _8bit-width_ register
#define MASK(a, b) (((uint8_t)-1 >> (7-(b))) & ~((1U<<(a))-1))

// Registers of simpleChip

#define REG_CONTROL 0x0 // 0
	#define REG_BYPASS_BIT BIT(1)
	#define REG_BYPASS_POS 1
	#define REG_BYPASS_VAL(rv) (((rv) & BIT(1)) >> 1)
		#define REG_BYPASS_ENABLE	1	// 0b1	0x1
		#define REG_BYPASS_DISABLE	0	// 0b0	0x0
	#define REG_FREE_RUN_STR 6
	#define REG_FREE_RUN_END 7
	#define REG_FREE_RUN_MSK MASK(6, 7)
	#define REG_FREE_RUN_VAL(rv) (((rv) & REG_FREE_RUN_MSK) >> 6)
	#define REG_FREE_RUN_SFT(v) (((v) & MASK(0, 1)) << 6)

#define REG_2 0x2 // 2
	#define REG_TYPE_STR 4
	#define REG_TYPE_END 5
	#define REG_TYPE_MSK MASK(4, 5)
	#define REG_TYPE_VAL(rv) (((rv) & REG_TYPE_MSK) >> 4)
	#define REG_TYPE_SFT(v) (((v) & MASK(0, 1)) << 4)
		#define REG_TYPE_CLIENT	0	// 0b0	0x0
		#define REG_TYPE_SERVER	1	// 0b1	0x1
		#define REG_TYPE_ROUTE	2	// 0b10	0x2
		#define REG_TYPE_PEER	3	// 0b11	0x3
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
  field, including name, offset and enumrations (optional) in the format of
  `name: offset (enumVal: enumName,)` that contains two parts. The `name-offset`
  part is mandatory contrasting to the `enums` part. Below explains more details
  about `FIELD` line.
  
  Example: `ck_prior1 : 0-1`, where the field's name is `ck_prior1`, offset
  range is `0-1` and there is no `enums` part.
  
## `name-offset` part

This part of `FIELD` line contains contains field's name and offset. The offset
are either bit or range which use `-` to connect. The start bit and end bit can
in either sides of `-`.
  
## `enums` part

This part of `FIELD` line is optional and gives the enumeration range this
field can have, which is usually used when the field's offset is not only one
bit.

Embraced by `(` and `)`, the enums are paired with `:`, where the enum value 
is at left and value name right. Enum pairs are separated with `,`.

For example, the chip spec have one field like:

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

As the example shows, the enum number can be decimal and binary. Actually, it
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
field_1: bit
field_2: startBit - endBit
field_3: endBit - startBit (val1: val1Name, val2: val2Name)
```

# output format

Right now, the only one output format is the C format which uses `#define`
macros to represent all registers and fields and capitalize all their names.

It handle two types of field: `bit` and `range`. There are some common macros and
some are not. The common macros are:

- `REG_XXX_VAL(rv)`: For `bit` type, it gets the value of this field, 1 or 0 of
  course. For `range`, it can get this field's value from register value `rv`.
  For example, `REG_XXX_VAL(0x10) = 2` for `3-4`.

- `REG_XXX_VALNAME VALNUM`: It is often used for `range` type. As the above
  example `type: 4-5 (0b00: client, 0x01: server, 2: route, 3: peer)`, the ouput
  is:
  ```
		#define REG_TYPE_CLIENT	0	// 0b0	0x0
		#define REG_TYPE_SERVER	1	// 0b1	0x1
		#define REG_TYPE_ROUTE	2	// 0b10	0x2
		#define REG_TYPE_PEER	3	// 0b11	0x3
  ```

- The macros only for `bit` type are:
	- `REG_XXX_BIT`: i.e. BIT(n), n is the offset.
	- `REG_XXX_POS`: it is the offset of this field.

- The macros only for `range` type are:
	- `REG_XXX_MSK`: the mask for this field. For example, `REG_XXX_MSK = 0x18`
	  for `3-4`.
	- `REG_XXX_SFT(v)`: it shifts the field's value `v` to the correct offset.
	  For example, `REG_XXX_SFT(1) = 0x08` for `3-4`.
