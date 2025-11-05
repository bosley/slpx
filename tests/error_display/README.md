# Error Display Test Cases

This directory contains test files demonstrating the precise error location reporting in SLPX.

## Test Files

### `basic.slpx`
**Error:** Wrong argument count to `set`
- Shows error on line 3, column 6
- Points to the first extra argument

### `argcount.slpx`  
**Error:** Wrong argument count to `int/add`
- Shows insufficient arguments to arithmetic function

### `typeerror.slpx`
**Error:** Type mismatch in `int/add`
- Shows attempting to add string to integer
- Points to the specific string argument that caused the type error

### `undefined.slpx`
**Error:** Undefined identifier
- Shows reference to `undefined_var` on line 3, column 19
- Points exactly to where the undefined variable is used

### `notfound.slpx`
**Error:** Function not found
- Shows call to undefined function on line 3, column 2
- Points to the function name that doesn't exist

### `badcall.slpx`
**Error:** Undefined identifier (non-existent function)
- Similar to notfound.slpx but with arguments
- Shows line 4, column 2

### `comprehensive.slpx`
**Error:** Wrong argument count to user-defined function
- Tests that position tracking works for function definitions
- Shows error at line 5, column 4

### `unclosed.slpx`
**Error:** Parse error - unclosed list
- Tests that parse errors (not runtime errors) are properly formatted
- Shows error at line 3, column 1 pointing to the opening parenthesis

## Error Types

### Parse Errors
Caught during parsing (syntax errors):
- Unclosed lists/strings
- Invalid syntax
- Macro errors

### Runtime Errors  
Caught during evaluation (semantic errors):
- Type mismatches
- Undefined identifiers
- Wrong argument counts

## Error Format

All errors display:
1. **File path** - Full path to the source file
2. **Line number** - The line where the error occurred
3. **Column number** - The exact column position
4. **Source line** - The actual line of code
5. **Pointer** - A `^` character pointing to the error location
6. **Error message** - Clear description of what went wrong

Example:
```
Error in /path/to/file.slpx at line 3, column 6:
  3 | (set a 3 3) ; fail, requires two arguments
           ^
wrong number of arguments: expected 2, got 3
```

