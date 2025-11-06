# CGS Bits Functions (`bits`)

Bit-level manipulation functions for SLPX. Provides conversion between numeric types and their bit representations.

## Function Reference

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `bits/explode` | `value :*` | `:L` | Converts an integer or real to a list of 64 bits (0 or 1). Returns error for unsupported types. |
| `bits/int` | `bits :L` | `:I` | Converts a list of 64 bits to a signed 64-bit integer. Returns error if list is not 64 elements or contains non-integer/non-binary values. |
| `bits/real` | `bits :L` | `:R` | Converts a list of 64 bits to a 64-bit floating point number. Returns error if list is not 64 elements or contains non-integer/non-binary values. |

## Type Legend

- `:*` - Any type
- `:I` - Integer (64-bit signed)
- `:R` - Real (64-bit floating point)
- `:L` - List
- `:E` - Error

## Notes

### Bit Representation

All conversions use 64-bit representations:
- Integers are treated as unsigned 64-bit values during bit extraction
- Reals use IEEE 754 double-precision format
- Bits are ordered from LSB (index 0) to MSB (index 63)

### Return Values

Functions return:
- List of 64 integers (0 or 1) for `bits/explode`
- Integer value for `bits/int`
- Real value for `bits/real`
- Error object for invalid inputs or conversions

### Error Conditions

- `bits/explode`: Returns error if value is not integer or real type
- `bits/int`: Returns error if list length is not 64, contains non-integers, or contains values other than 0 or 1
- `bits/real`: Returns error if list length is not 64, contains non-integers, or contains values other than 0 or 1

