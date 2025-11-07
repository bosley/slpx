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

### Random Number Generation

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `int/rand` | `lower :I`, `upper :I` | `:I` | Generate random integer in range [lower, upper] (inclusive). Returns error if lower > upper. |
| `real/rand` | `lower :R`, `upper :R` | `:R` | Generate random real number in range [lower, upper). Returns error if lower > upper. |

### Advanced Math Functions

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `real/sqrt` | `value :R` | `:R` | Square root. Returns error on negative input. |
| `real/exp` | `value :R` | `:R` | Exponential function (e^x). Returns error on overflow. |
| `real/log` | `value :R` | `:R` | Natural logarithm (ln). Returns error on non-positive input. |
| `real/ceil` | `value :R` | `:I` | Ceiling function - round up to nearest integer. |
| `real/round` | `value :R` | `:I` | Round to nearest integer (half away from zero). |

### Real Number Inspection

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `real/is-nan` | `value :R` | `:I` | Check if value is NaN. Returns `1` if NaN, `0` otherwise. |
| `real/is-inf` | `value :R` | `:I` | Check if value is infinite. Returns `1` if infinite, `0` otherwise. |
| `real/is-finite` | `value :R` | `:I` | Check if value is finite (not NaN or infinite). Returns `1` if finite, `0` otherwise. |

### Absolute Value

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `int/abs` | `value :I` | `:I` | Absolute value of integer. |
| `real/abs` | `value :R` | `:R` | Absolute value of real number. |

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

**Random Number Generation:**
- `int/rand` - Lower bound greater than upper bound
- `real/rand` - Lower bound greater than upper bound

**Advanced Math Functions:**
- `real/sqrt` - Negative input
- `real/exp` - Result overflow (infinity)
- `real/log` - Zero or negative input

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

### Random Number Generation
```lisp
(int/rand 1 10)                   ; Random integer between 1 and 10 (inclusive)
(int/rand 42 100)                 ; Random integer between 42 and 100 (inclusive)
(int/rand 5 5)                    ; Always returns 5

(set lower 42)
(set r (int/rand lower 100))      ; Random integer between 42 and 100

(real/rand 0.0 1.0)               ; Random real number between 0.0 and 1.0
(real/rand 3.33 100.0)            ; Random real number between 3.33 and 100.0
(real/rand 2.5 2.5)               ; Always returns 2.5

(set real_lower 3.33)
(set real_upper 100.0)
(set r (real/rand real_lower real_upper))  ; Random real between bounds
```

### Advanced Math Functions
```lisp
(real/sqrt 4.0)                   ; 2.0
(real/sqrt 2.0)                   ; 1.414...
(real/sqrt 0.0)                   ; 0.0

(real/exp 0.0)                    ; 1.0 (e^0)
(real/exp 1.0)                    ; 2.718... (e^1)
(real/exp 2.0)                    ; 7.389... (e^2)

(real/log 1.0)                    ; 0.0 (ln(1))
(real/log 2.718281828459045)      ; 1.0 (ln(e))
(real/log 10.0)                   ; 2.302... (ln(10))

(real/ceil 3.14)                  ; 4
(real/ceil -3.14)                 ; -3 (rounds toward positive infinity)
(real/ceil 5.0)                   ; 5

(real/round 3.5)                  ; 4 (half away from zero)
(real/round 3.4)                  ; 3
(real/round -3.5)                 ; -4 (half away from zero)
```

### Absolute Value
```lisp
(int/abs 42)                      ; 42
(int/abs -42)                     ; 42
(int/abs 0)                       ; 0

(real/abs 3.14)                   ; 3.14
(real/abs -3.14)                  ; 3.14
(real/abs 0.0)                    ; 0.0
```

### Real Number Inspection
```lisp
(real/is-finite 3.14)             ; 1 (true)
(real/is-finite 0.0)              ; 1 (true)

(set result (real/div 1.0 0.0))   ; Creates infinity (if not caught)
(real/is-inf result)              ; 1 (true)
(real/is-finite result)           ; 0 (false)

(set ratio (real/div 0.0 0.0))    ; Creates NaN (if not caught)
(real/is-nan ratio)               ; 1 (true)
(real/is-finite ratio)            ; 0 (false)
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

(real/sqrt (int/real (int/abs -16)))      ; sqrt(abs(-16)) = 4.0

(real/exp (real/log 10.0))                ; e^(ln(10)) = 10.0

(real/round (real/sqrt 50.0))             ; round(sqrt(50)) = 7

(int/abs (int/sub 
  (real/ceil 3.7) 
  (real/round 8.5)))                      ; abs(4 - 9) = 5
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

**Random Number Generation:**
- `int/rand` generates random integers in the inclusive range [lower, upper]
- `real/rand` generates random real numbers in the range [lower, upper)
- Both functions return the bound value directly if lower equals upper
- Uses `math/rand/v2` for high-quality pseudo-random number generation
- Integer random uses uniform distribution via `rand.IntN`
- Real random uses uniform distribution via `rand.Float64`

**Advanced Math Functions:**
- `real/sqrt` uses `math.Sqrt` for square root calculation
- `real/exp` uses `math.Exp` for exponential calculation, checks for overflow
- `real/log` uses `math.Log` for natural logarithm
- `real/ceil` uses `math.Ceil` and converts to integer (rounds toward positive infinity)
- `real/round` uses `math.Round` and converts to integer (ties away from zero)

**Real Number Inspection:**
- `real/is-nan` uses `math.IsNaN` to detect IEEE 754 NaN values
- `real/is-inf` uses `math.IsInf` to detect positive or negative infinity
- `real/is-finite` returns true only if value is not NaN and not infinite
- All inspection functions return integer boolean (1 for true, 0 for false)

**Absolute Value:**
- `int/abs` uses conditional check and negation for integers
- `real/abs` uses `math.Abs` for real numbers
- Both preserve the magnitude while removing sign

