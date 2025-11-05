# CGS Numbers Functions (`arith`)

Arithmetic, comparison, and type conversion command group for integers and real numbers in SLPX.

## Function Reference

### Integer Arithmetic

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `int/add` | `a :I`, `b :I` | `:I` | Add two integers. |
| `int/sub` | `a :I`, `b :I` | `:I` | Subtract b from a. |
| `int/mul` | `a :I`, `b :I` | `:I` | Multiply two integers. |
| `int/div` | `a :I`, `b :I` | `:I` | Divide a by b (integer division). Returns error on division by zero. |
| `int/mod` | `a :I`, `b :I` | `:I` | Modulo operation (a mod b). Returns error on modulo by zero. |
| `int/pow` | `a :I`, `b :I` | `:I` | Raise a to the power of b. Returns error on negative exponent. |
| `int/sum` | `values :I...` | `:I` | Sum multiple integers (variadic). Requires at least one argument. |

### Real Arithmetic

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `real/add` | `a :R`, `b :R` | `:R` | Add two real numbers. |
| `real/sub` | `a :R`, `b :R` | `:R` | Subtract b from a. |
| `real/mul` | `a :R`, `b :R` | `:R` | Multiply two real numbers. |
| `real/div` | `a :R`, `b :R` | `:R` | Divide a by b. Returns error on division by zero. |
| `real/pow` | `a :R`, `b :R` | `:R` | Raise a to the power of b. Returns error on NaN or Inf result. |
| `real/sum` | `values :R...` | `:R` | Sum multiple real numbers (variadic). Requires at least one argument. |

### Type Conversions

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `int/real` | `value :I` | `:R` | Convert integer to real number. |
| `real/int` | `value :R` | `:I` | Convert real to integer. Floors the value before conversion. |

### Integer Comparisons

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `int/eq` | `a :I`, `b :I` | `:I` | Equality comparison. Returns `1` if equal, `0` if not. |
| `int/gt` | `a :I`, `b :I` | `:I` | Greater than comparison. Returns `1` if a > b, `0` otherwise. |
| `int/gte` | `a :I`, `b :I` | `:I` | Greater than or equal comparison. Returns `1` if a >= b, `0` otherwise. |
| `int/lt` | `a :I`, `b :I` | `:I` | Less than comparison. Returns `1` if a < b, `0` otherwise. |
| `int/lte` | `a :I`, `b :I` | `:I` | Less than or equal comparison. Returns `1` if a <= b, `0` otherwise. |

### Real Comparisons

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `real/eq` | `a :R`, `b :R` | `:I` | Equality comparison. Returns `1` if equal, `0` if not. |
| `real/gt` | `a :R`, `b :R` | `:I` | Greater than comparison. Returns `1` if a > b, `0` otherwise. |
| `real/gte` | `a :R`, `b :R` | `:I` | Greater than or equal comparison. Returns `1` if a >= b, `0` otherwise. |
| `real/lt` | `a :R`, `b :R` | `:I` | Less than comparison. Returns `1` if a < b, `0` otherwise. |
| `real/lte` | `a :R`, `b :R` | `:I` | Less than or equal comparison. Returns `1` if a <= b, `0` otherwise. |

## Type Legend

- `:I` - Integer (64-bit signed)
- `:R` - Real (64-bit floating-point)
- `...` - Variadic (accepts multiple arguments)

## Notes

### Boolean Representation

All comparison functions return integers as boolean values:
- `1` = true
- `0` = false

This allows comparisons to be used directly in conditional logic.

### Error Handling

Functions that can fail return error objects:

**Integer Operations:**
- `int/div` - Division by zero
- `int/mod` - Modulo by zero
- `int/pow` - Negative exponent (use `real/pow` for negative exponents)

**Real Operations:**
- `real/div` - Division by zero
- `real/pow` - Result is NaN or Infinity

### Type Conversions

**`real/int` Flooring Behavior:**
- Uses `math.Floor` before conversion
- `3.14` → `3`
- `-7.8` → `-8` (floors towards negative infinity)

**`int/real` Conversion:**
- Lossless conversion for values within float64 precision range

### Variadic Functions

**`int/sum` and `real/sum`:**
- Accept one or more arguments
- Validate all arguments are of correct type at runtime
- Type mismatch returns error with position information

### Integer Power Implementation

`int/pow` uses fast exponentiation by squaring (binary exponentiation):
- Efficient for large exponents
- Only supports non-negative exponents
- For negative exponents, use: `(real/pow (int/real a) (int/real b))`

## Examples

### Basic Arithmetic
```lisp
(int/add 10 5)                    ; 15
(int/sub 10 5)                    ; 5
(int/mul 10 5)                    ; 50
(int/div 10 3)                    ; 3 (integer division)
(int/mod 10 3)                    ; 1
(int/pow 2 10)                    ; 1024

(real/add 10.5 5.2)               ; 15.7
(real/sub 10.5 5.2)               ; 5.3
(real/mul 10.5 2.0)               ; 21.0
(real/div 10.0 3.0)               ; 3.333...
(real/pow 2.0 0.5)                ; 1.414... (square root)
```

### Variadic Operations
```lisp
(int/sum 1 2 3 4 5)               ; 15
(int/sum 10 20 30)                ; 60
(int/sum 42)                      ; 42

(real/sum 1.5 2.5 3.5)            ; 7.5
(real/sum 10.0 20.0 30.0)         ; 60.0
(real/sum 5.5)                    ; 5.5
```

### Type Conversions
```lisp
(int/real 42)                     ; 42.0
(real/int 3.14)                   ; 3
(real/int -7.8)                   ; -8
(real/int (real/div 7.0 2.0))     ; 3 (7.0/2.0 = 3.5, floored to 3)
```

### Comparisons
```lisp
(int/eq 10 10)                    ; 1 (true)
(int/eq 10 20)                    ; 0 (false)
(int/gt 20 10)                    ; 1 (true)
(int/lt 10 20)                    ; 1 (true)
(int/gte 10 10)                   ; 1 (true)
(int/lte 10 10)                   ; 1 (true)

(real/eq 3.14 3.14)               ; 1 (true)
(real/gt 3.14 2.71)               ; 1 (true)
(real/lt 2.71 3.14)               ; 1 (true)
```

### Complex Expressions
```lisp
(int/mul (int/add 5 3) (int/sub 10 2))    ; (5+3) * (10-2) = 64

(int/sum 
  (int/pow 1 2) 
  (int/pow 2 2) 
  (int/pow 3 2) 
  (int/pow 4 2))                          ; 1+4+9+16 = 30

(real/int (real/div 
  (real/add 10.5 5.5) 
  (real/add 2.0 2.0)))                    ; (10.5+5.5)/(2.0+2.0) = 4.0 → 4
```

### Error Handling with Comparisons
```lisp
(set is_positive (int/gt x 0))
(if (int/eq is_positive 1)
  (putln "positive")
  (putln "not positive"))
```

### Safe Division
```lisp
(try 
  (int/div 10 0)
  (putln "Division by zero error caught"))
```

## Performance Notes

- Integer operations are exact (no floating-point errors)
- `int/pow` uses O(log n) algorithm for efficiency
- Real operations may have floating-point precision limitations
- Comparison functions are optimized single operations
- Variadic functions validate types at runtime (small overhead)

## Implementation Details

**Variadic Type Safety:**
Both `int/sum` and `real/sum` validate each argument's type at runtime and return descriptive errors:
```
"int/sum: all arguments must be integers, got <type> at position <N>"
```

**Integer Division:**
Uses Go's integer division (truncation towards zero):
- `7 / 2 = 3`
- `-7 / 2 = -3`

**Real Floor Conversion:**
Uses `math.Floor` (rounds towards negative infinity):
- `3.9 → 3`
- `-3.1 → -4`

