# .svd input format

Accept .svd as input format

# struct output

Add `-f` format option to specify output format. Currently, the `cmacro` is 
supported. Next, `cstruct` should be added to output structure of a register,
like:

And access a bit-field of a register, we can use structure instead.
For example, with ARM 32-bit, we can use

```c
/* Description of register */
typedef union {
	struct {
	UINT32 EN: 1; /* Description of bit field */
	UINT32 Reserved: 31; /* Reserved bit field */
	} BIT;
	UINT32 INT;
} drv<module_name>_<register_name>_t;
```
```c
/* ================ TIMER0 ================ */
typedef struct {                                    
  __IO uint32_t  CR;                                
  __IO uint16_t  SR;                                
  __I  uint16_t  RESERVED0[5];
  __IO uint16_t  INT;                               
  __I  uint16_t  RESERVED1[7];
  __IO uint32_t  COUNT;                             
  __IO uint32_t  MATCH;                             
  union {
    __O  uint32_t  PRESCALE_WR;                     
    __I  uint32_t  PRESCALE_RD;                     
  };
  __I  uint32_t  RESERVED2[9];
  __IO uint32_t  RELOAD[4];                         
} TIMER0_Type;
```
# line number

Store line number and show them when the .regs format is invalid.

# DONE

- bit width
- license banner
- file flag
- module
- full-length macro
