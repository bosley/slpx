# CGS String Functions (`str`)

String manipulation and conversion command group for SLPX.

## Function Reference

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `str/eq` | `a :S`, `b :S` | `:I` | Compare two strings for equality. Returns `1` if equal, `0` if not. |
| `str/len` | `s :S` | `:I` | Get the length of a string (counts runes/characters). |
| `str/clear` | `s :S` | `:S` | Returns an empty string (ignores input). |
| `str/from` | `obj :*` | `:S` | Convert any object to its string representation. Uses precision setting for real numbers. |
| `str/int` | `s :S` | `:I` | Parse a string to an integer. Returns error if parsing fails. |
| `str/real` | `s :S` | `:R` | Parse a string to a real number. Returns error if parsing fails. |
| `str/list` | `s :S` | `:L` | Convert string to a list of individual character strings. |
| `str/concat` | `strings :S...` | `:S` | Concatenate multiple strings together (variadic). |
| `str/upper` | `s :S` | `:S` | Convert string to uppercase. |
| `str/lower` | `s :S` | `:S` | Convert string to lowercase. |
| `str/trim` | `s :S` | `:S` | Remove leading and trailing whitespace. |
| `str/contains` | `s :S`, `substr :S` | `:I` | Check if string contains substring. Returns `1` if found, `0` if not. |
| `str/index` | `s :S`, `substr :S` | `:I` | Find the byte index of first occurrence of substring. Returns `-1` if not found. |
| `str/slice` | `s :S`, `start :I`, `end :I` | `:S` | Extract substring from start to end index (rune-based, bounds-safe). |
| `str/split` | `s :S`, `sep :S` | `:L` | Split string by separator into a list of strings. |
| `str/replace` | `s :S`, `old :S`, `new :S` | `:S` | Replace all occurrences of old substring with new substring. |
| `str/precision` | `p :I` | `:I` | Set floating-point precision (0-255) for `str/from`. Returns the set precision. |

## Type Legend

- `:S` - String
- `:I` - Integer
- `:R` - Real (floating-point)
- `:L` - List
- `:*` - Any type
- `...` - Variadic (accepts multiple arguments)

## Notes

### Precision Control

The `str/precision` function controls how many decimal places are used when converting real numbers to strings via `str/from`:

```lisp
(str/precision 2)
(str/from 3.14159)  ; Returns "3.14"

(str/precision 10)
(str/from 3.14159)  ; Returns "3.1415900000"
```

Default precision is `6`. The precision setting is:
- Thread-safe (uses mutex internally)
- Persistent across multiple calls
- Bounded to range 0-255

### Unicode Support

The following functions properly handle Unicode/multi-byte characters by operating on runes:
- `str/len` - counts characters, not bytes
- `str/slice` - slices by character position, not byte position
- `str/list` - splits into individual characters correctly
- `str/upper` - correctly uppercases unicode characters
- `str/lower` - correctly lowercases unicode characters

**Note:** `str/index` returns byte positions, not rune positions. For ASCII strings this distinction doesn't matter, but for strings with multi-byte unicode characters, the byte index may differ from the character index.

### Error Handling

Functions that can fail return error objects:
- `str/int` - Returns error if string cannot be parsed as integer
- `str/real` - Returns error if string cannot be parsed as real number

### Variadic Functions

- `str/concat` - Accepts zero or more string arguments. Returns empty string if no arguments provided. Validates all arguments are strings at runtime.

## Examples

### Basic Operations
```lisp
(str/eq "hello" "hello")           ; 1
(str/len "Hello")                  ; 5
(str/upper "hello")                ; "HELLO"
(str/lower "WORLD")                ; "world"
```

### String Building
```lisp
(str/concat "Hello" " " "World")   ; "Hello World"
(str/from 42)                      ; "42"
(str/from 3.14)                    ; "3.140000" (with default precision)
```

### Parsing
```lisp
(str/int "42")                     ; 42
(str/real "3.14")                  ; 3.14
(str/list "Hi")                    ; ("H" "i")
```

### Search & Manipulation
```lisp
(str/contains "hello world" "world")      ; 1
(str/index "hello world" "world")         ; 6
(str/slice "hello world" 0 5)             ; "hello"
(str/split "a,b,c" ",")                   ; ("a" "b" "c")
(str/replace "hello world" "world" "there") ; "hello there"
```

### Advanced Usage
```lisp
(set name "   john doe   ")
(set cleaned (str/trim name))
(set upper (str/upper cleaned))
(set parts (str/split upper " "))
; parts = ("JOHN" "DOE")
```

