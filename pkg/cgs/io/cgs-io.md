# CGS IO Functions (`io`)

Input/output operations, color formatting, and console interaction command group for SLPX.

## Function Reference

### Output Operations

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `io/out` | `args :*...` | `:N` | Write arguments to output. Variadic function that converts all arguments to strings and flushes after each. |
| `io/out/set_precision` | `precision :I` | `:N` | Set decimal precision for real number output (0-20, default 6). Values outside range are clamped. |
| `io/flush` | - | `:N` | Flush output buffer. Returns error if flush fails. |

### Input Operations

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `io/in` | `prompt :S` | `:S` | Display prompt and read line of text input. Returns error if read fails. |
| `io/in/int` | `prompt :S` | `:I` | Display prompt and read integer input. Returns error if input is not a valid integer. |
| `io/in/real` | `prompt :S` | `:R` | Display prompt and read real number input. Returns error if input is not a valid number. |

### Color Operations

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `io/color/fg` | `color :S` | `:S` | Generate ANSI escape sequence for foreground color from hex string. Returns error if format invalid. |
| `io/color/bg` | `color :S` | `:S` | Generate ANSI escape sequence for background color from hex string. Returns error if format invalid. |
| `io/color/reset` | - | `:S` | Generate ANSI escape sequence to reset all colors and formatting. |

## Type Legend

- `:S` - String
- `:I` - Integer
- `:R` - Real (floating-point number)
- `:N` - None (no return value)
- `:*` - Any type
- `...` - Variadic (accepts multiple arguments)

## Notes

### Output Behavior

**`io/out` Type Conversion:**
The output function handles different types automatically:
- **String:** Output as-is
- **Integer:** Formatted as decimal number (e.g., `42`)
- **Real:** Formatted with precision setting (e.g., `3.141593` with precision 6)
- **Other types:** Converted using `Encode()` method

**Auto-Flush:**
`io/out` automatically flushes after each argument to ensure immediate output visibility.

**Precision Control:**
- Default precision: 6
- Valid range: 0-20 (values outside are automatically clamped)
- Affects all subsequent `io/out` calls with real numbers

### Input Behavior

**Blocking Operations:**
All input functions block until user provides input and presses Enter.

**String Input:**
`io/in` returns the complete line including spaces, with newline removed.

**Numeric Input Parsing:**
- `io/in/int`: Parses base-10 integers using `strconv.ParseInt(line, 10, 64)`
- `io/in/real`: Parses floating-point numbers using `strconv.ParseFloat(line, 64)`
- Leading/trailing whitespace is automatically trimmed before parsing

### Color Formatting

**Hex Color Format:**
Color functions accept hex color strings in format `#RRGGBB` or `RRGGBB`:
- Exactly 6 hexadecimal digits required
- Hash prefix is optional and automatically stripped
- Returns error object if format is invalid

**ANSI Escape Sequences:**
Color functions return ANSI escape code strings:
- `io/color/fg`: `\033[38;2;R;G;Bm` (24-bit foreground color)
- `io/color/bg`: `\033[48;2;R;G;Bm` (24-bit background color)
- `io/color/reset`: `\033[0m` (reset all attributes)

**Terminal Compatibility:**
Requires terminal with 24-bit true color support. Older terminals may ignore or misinterpret escape sequences.

### Error Handling

Functions that can fail return error objects with `Type: OBJ_TYPE_ERROR`:
- `io/in` - Read failures
- `io/in/int` - Read failures or input not valid integer format
- `io/in/real` - Read failures or input not valid number format
- `io/color/fg` - Invalid hex color format (not 6 characters or invalid hex digits)
- `io/color/bg` - Invalid hex color format (not 6 characters or invalid hex digits)
- `io/flush` - Output buffer flush failures
- `io/out/set_precision` - Never fails (clamps to valid range)

Error messages include descriptive text and position 0.

## Examples

### Basic Output

```lisp
(io/out "Hello, World!")
```

Output: `Hello, World!`

```lisp
(io/out "The answer is: " 42)
```

Output: `The answer is: 42`

```lisp
(io/out "Multiple " "arguments " "work")
```

Output: `Multiple arguments work`

### Integer Output

```lisp
(io/out 123)
```

Output: `123`

```lisp
(io/out -456)
```

Output: `-456`

### Real Number Output

```lisp
(io/out 3.14159265359)
```

Output: `3.141593` (with default precision 6)

### Precision Control

```lisp
(io/out/set_precision 2)
(io/out 3.14159)
```

Output: `3.14`

```lisp
(io/out/set_precision 8)
(io/out 3.14159265359)
```

Output: `3.14159265`

```lisp
(io/out/set_precision 0)
(io/out 3.14159)
```

Output: `3`

```lisp
(io/out/set_precision 25)
(io/out 1.5)
```

Output: `1.50000000000000000000` (clamped to precision 20)

### String Input

```lisp
(io/in "Enter your name: ")
```

Displays prompt and waits for input. Returns entered string.

### Integer Input

```lisp
(io/in/int "Enter age: ")
```

Displays prompt and waits for integer. Returns integer or error object.

Valid input: `25` → Returns `25` as integer
Invalid input: `abc` → Returns error: "input is not a valid integer: abc"

### Real Number Input

```lisp
(io/in/real "Enter temperature: ")
```

Displays prompt and waits for number. Returns real or error object.

Valid input: `98.6` → Returns `98.6` as real
Invalid input: `hot` → Returns error: "input is not a valid real number: hot"

### Foreground Colors

```lisp
(io/out (io/color/fg "#FF0000") "Red text" (io/color/reset))
```

Output: `Red text` (in red, then reset)

```lisp
(io/out (io/color/fg "00FF00") "Green text" (io/color/reset))
```

Output: `Green text` (in green, hash optional)

```lisp
(io/out (io/color/fg "#0000FF") "Blue" (io/color/reset))
```

Output: `Blue` (in blue)

### Background Colors

```lisp
(io/out (io/color/bg "#FF0000") "Red background" (io/color/reset))
```

Output: `Red background` (with red background)

```lisp
(io/out (io/color/fg "#FFFFFF") (io/color/bg "#000000") "White on black" (io/color/reset))
```

Output: `White on black` (white text on black background)

### Color Reset

```lisp
(io/out (io/color/fg "#FF0000") "Red" (io/color/reset) " Normal")
```

Output: `Red Normal` (red, then default color)

```lisp
(io/color/reset)
```

Returns: `"\033[0m"`

### Invalid Color Formats

```lisp
(io/color/fg "#FF00GG")
```

Returns error: "invalid hex color: invalid syntax"

```lisp
(io/color/fg "#F00")
```

Returns error: "invalid hex color: hex color must be 6 characters (got 3)"

```lisp
(io/color/bg "GGGGGG")
```

Returns error: "invalid hex color: invalid syntax"

### Flush

```lisp
(io/flush)
```

Flushes output buffer. Returns none or error object if flush fails.

### Combined Example

```lisp
(io/out (io/color/fg "#00FFFF") "Cyan text " (io/color/reset) "Normal text")
```

Output: `Cyan text Normal text` (first part cyan, second part normal)

```lisp
(io/out/set_precision 3)
(io/out "Pi is approximately " 3.14159265 "\n")
(io/flush)
```

Output: `Pi is approximately 3.142` (with newline, then flushed)

## Performance Notes

- `io/out` - O(1) per argument, includes automatic flush
- `io/in` - Blocking I/O operation
- `io/in/int` - Blocking I/O + O(n) parsing where n is input length
- `io/in/real` - Blocking I/O + O(n) parsing where n is input length
- `io/flush` - O(1) buffer flush operation
- `io/color/fg` - O(1) hex parsing and string formatting
- `io/color/bg` - O(1) hex parsing and string formatting
- `io/color/reset` - O(1) constant string return
- `io/out/set_precision` - O(1) state update

## Implementation Details

**Output Buffering:**
The IO system uses a buffered writer accessed via `runtime.GetIO()`. `io/out` calls `Flush()` automatically after each argument write.

**Precision State:**
The precision setting is stored in the `ioFunctions` struct instance and persists across calls. Stored as integer, clamped between 0-20 inclusive.

**Input Trimming:**
Numeric input functions (`io/in/int`, `io/in/real`) use `strings.TrimSpace()` before parsing. String input (`io/in`) returns line verbatim minus newline.

**Color Parsing:**
Hex colors parsed via:
1. Strip `#` prefix with `strings.TrimPrefix()`
2. Validate length == 6
3. Parse as unsigned 32-bit integer: `strconv.ParseUint(hexColor, 16, 32)`
4. Extract RGB: `r = (val >> 16) & 0xFF`, `g = (val >> 8) & 0xFF`, `b = val & 0xFF`
5. Format into ANSI escape sequence

**Error Return Pattern:**
Functions return `(object.Obj, error)` where the error is always `nil` and errors are encoded as `object.Obj` with `Type: OBJ_TYPE_ERROR`, `Position: 0`, and descriptive message string.

**IO Interface:**
All functions interact with `env.IO` interface obtained from runtime context, providing methods:
- `WriteString(string)` - Write string to output
- `ReadLine() (string, error)` - Read line from input
- `Flush() error` - Flush output buffer
