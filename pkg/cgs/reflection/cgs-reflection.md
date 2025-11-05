# CGS Reflection Functions (`reflect`)

Type introspection and runtime type checking command group for SLPX.

## Function Reference

### Type Inspection

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `reflect/type?` | `value :*` | `:S` | Returns the type name of a value as a string. Does not evaluate the argument. |
| `reflect/equal?` | `a :*`, `b :*` | `:I` | Returns `1` if both values have the same type, `0` otherwise. |

### Type Predicates

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `reflect/int?` | `value :*` | `:I` | Returns `1` if value is an integer, `0` otherwise. |
| `reflect/real?` | `value :*` | `:I` | Returns `1` if value is a real number, `0` otherwise. |
| `reflect/str?` | `value :*` | `:I` | Returns `1` if value is a string, `0` otherwise. |
| `reflect/list?` | `value :*` | `:I` | Returns `1` if value is a list, `0` otherwise. |
| `reflect/fn?` | `value :*` | `:I` | Returns `1` if value is a function, `0` otherwise. |
| `reflect/none?` | `value :*` | `:I` | Returns `1` if value is none, `0` otherwise. |
| `reflect/error?` | `value :*` | `:I` | Returns `1` if value is an error, `0` otherwise. |
| `reflect/some?` | `value :*` | `:I` | Returns `1` if value is quoted (some), `0` otherwise. |
| `reflect/ident?` | `value :*` | `:I` | Returns `1` if value is an identifier, `0` otherwise. |

## Type Legend

- `:*` - Any type (no type checking)
- `:I` - Integer (64-bit signed)
- `:S` - String
- Other types: `:R` (Real), `:L` (List), `:F` (Function), `:_` (None), `:E` (Error), `:Q` (Some/Quoted)

## Notes

### Boolean Representation

All type predicate functions return integers as boolean values:
- `1` = true (type matches)
- `0` = false (type does not match)

This allows predicates to be used directly in conditional logic.

### Type Name Strings

`reflect/type?` returns type names as strings:
- `"integer"` - Integer values
- `"real"` - Real (floating-point) values
- `"string"` - String values
- `"list"` - List values
- `"function"` - Function values
- `"none"` - None values
- `"error"` - Error values
- `"some"` - Quoted (some) values
- `"identifier"` - Identifier values

### Special Evaluation Behavior

**`reflect/type?` does NOT evaluate its argument:**
- This allows inspecting the type of expressions without evaluating them
- Identifiers are NOT resolved automatically
- Lists are NOT executed
- To check the type of a variable's value, the argument will still be evaluated through normal identifier resolution

All other reflection functions evaluate their arguments normally.

### Type Equality

`reflect/equal?` compares only types, not values:
- `(reflect/equal? 42 100)` → `1` (both integers)
- `(reflect/equal? 42 "hello")` → `0` (different types)
- `(reflect/equal? '(1 2) '(3 4))` → `1` (both lists)

## Examples

### Basic Type Inspection

```lisp
(reflect/type? 42)                   ; "integer"
(reflect/type? 3.14)                 ; "real"
(reflect/type? "hello")              ; "string"
(reflect/type? '())                  ; "list"
(reflect/type? _)                    ; "none"
(reflect/type? (fn (..) :_))         ; "function"

(set x 100)
(reflect/type? x)                    ; "integer" (evaluates x to 100)
```

### Type Predicates

```lisp
(reflect/int? 42)                    ; 1 (true)
(reflect/int? 3.14)                  ; 0 (false)
(reflect/real? 3.14)                 ; 1 (true)
(reflect/str? "hello")               ; 1 (true)
(reflect/list? '(1 2 3))             ; 1 (true)
(reflect/fn? (fn (..) :_))           ; 1 (true)
(reflect/none? _)                    ; 1 (true)
(reflect/error? (uq '@(error msg)))  ; 1 (true)
(reflect/some? '(1 2))               ; 1 (true - quoted list)
```

### Type Equality Checks

```lisp
(set x 42)
(set y 100)
(set z "hello")

(reflect/equal? x y)                 ; 1 (both integers)
(reflect/equal? x z)                 ; 0 (integer vs string)
(reflect/equal? '(1 2) '(3 4))       ; 1 (both lists)
```

### Conditional Logic with Type Checks

```lisp
(set value 42)

(if (reflect/int? value)
  (putln "Value is an integer")
  (putln "Value is not an integer"))

(if (reflect/equal? value "text")
  (putln "Same type")
  (putln "Different types"))
```

### Generic Function with Type Checking

```lisp
(fn (process value) :*
  (if (reflect/int? value)
    (int/mul value 2)
    (if (reflect/str? value)
      (str/concat value " processed")
      _)))

(process 21)                         ; 42
(process "data")                     ; "data processed"
(process _)                          ; _
```

### Safe Type Coercion

```lisp
(fn (to_int value) :I
  (if (reflect/int? value)
    value
    (if (reflect/real? value)
      (real/int value)
      (if (reflect/str? value)
        (str/parse_int value)
        0))))
```

### Type-based Dispatch

```lisp
(fn (stringify value) :S
  (if (reflect/int? value)
    (int/str value)
    (if (reflect/real? value)
      (real/str value)
      (if (reflect/str? value)
        value
        (if (reflect/none? value)
          "none"
          "unknown")))))
```

### Validating Function Arguments

```lisp
(fn (safe_divide a b) :*
  (if (reflect/int? a)
    (if (reflect/int? b)
      (if (int/eq b 0)
        '@(division by zero)
        (int/div a b))
      '@(second argument must be integer))
    '@(first argument must be integer)))
```

### Checking for Errors

```lisp
(set result (some_operation))

(if (reflect/error? result)
  (putln "Operation failed")
  (putln "Operation succeeded"))
```

### Quoted Value Detection

```lisp
(reflect/some? '(1 2 3))             ; 1 (quoted list)
(reflect/some? (list 1 2 3))         ; 0 (evaluated list)
(reflect/some? 'x)                   ; 1 (quoted identifier)
(reflect/some? x)                    ; 0 (evaluated identifier)
```

## Performance Notes

- Type checking is a constant-time operation
- `reflect/type` performs minimal string allocation
- Type predicates are optimized for direct type comparison
- No evaluation overhead for `reflect/type` (arguments not evaluated)
- All other functions evaluate arguments once before type checking

## Implementation Details

**Evaluation Strategy:**
- `reflect/type?`: `EvaluateArgs: false` - receives raw AST node
- All other functions: `EvaluateArgs: true` - receives evaluated values

**Type Comparison:**
- Direct comparison of `object.ObjType` enum values
- No deep inspection of values (only type tags checked)

**Return Values:**
- Type predicates return `object.Integer(1)` or `object.Integer(0)`
- `reflect/type?` returns a `object.String` with type name
- `reflect/equal?` returns integer boolean (`1` or `0`)

