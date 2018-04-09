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
// Registers of si5324
#define BIT(x) (1 << (x))
#define MASK(a, b) (((uint8_t)-1 >> (7-(b))) & ~((1U<<(a))-1))

#define REG_CONTROL 0x0
        #define REG_BYPASS_REG_BIT BIT(1)
        #define REG_CKOUT_ALWAYS_ON_BIT BIT(5)
        #define REG_FREE_RUN_MSK MASK(6, 7)

#define REG_1 0x1
        #define REG_CK_PRIOR1_MSK MASK(0, 1)

#define REG_16 0x10
        #define REG_BWSEL_REG_MSK MASK(4, 7)
```

It contains two types of field: `bit` and `mask`.


# input format

The above example shows the input file's basic format which contains only a few
keywords and follows simple and loose rules. All keywords are case-insensitive.

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
