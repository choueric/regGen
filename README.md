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

It can handle two types of field: `bit` and `mask`:
- for `bit` type, it only generates one single macro `REG_XXX_BIT`
- for `mask` type, it outputs macros below:
	- `REG_XXX_MSK`: the mask for this field, for example `0x18` for `3-4`.
	- `REG_XXX_VAL(rv)`: this macro can get this field's value from register 
	  value `rv`. For example, `REG_XXX_VAL(0x10) = 2` for `3-4`.
	- `REG_XXX_SFT(v)`: it shifts the field's value `v` to the correct offset.
	  For example, `REG_XXX_SFT(3) = 0x18` for `3-4`.

# input format

The above example shows the input file's basic format which contains only a few
keywords and follows simple and loose rules. All keywords are case-insensitive.

- `#`: A line will be treated as comment line if containing `#`.
- `<chip>:si5324`: `<chip>` specify the chip's name and is optional
- `<REG>[Control]: 0`: `<reg>` start a description of a register and its fields
  until a new `<reg>`. The register's name is `Control` and is optional. In this
  case, the offset will be its name.
- `ck_prior1 : 0-1`, `CKOUT_ALWAYS_ON: 5`: The lines below `<reg>` define the
  fields in the format `name: offset`. It doesn't matter to indent or not.
  The offset are either bit or range which use `-` to connect. The start bit and
  end bit can in either sides of `-` due to the process of correction in program.
  
A more formal sepcification is like:
```
<chip>: {chip_name}

<reg>[{reg_name}]:{offset}
{filed_1}: {bit}
{filed_2}: {start_bit}-{end_bit}
{filed_3}: {end_bit}-{start_bit}
```

# output format

Right now, only output of C format is supported. These definitions are outputed
as macros, so the names of registers and fields are capitalized.
