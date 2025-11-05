# CGS List Functions (`list`)

List manipulation, iteration, and functional programming command group for SLPX.

## Function Reference

### Core Operations

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `list/new` | `length :I`, `default :*` | `:L` | Create new list of specified length filled with deep copies of default value. |
| `list/len` | `list :L` | `:I` | Get the length of a list. |
| `list/clear` | `list :L` | `:L` | Return an empty list. |
| `list/empty` | `list :L` | `:I` | Check if list is empty. Returns `1` if empty, `0` if not. |

### Element Access

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `list/get` | `list :L`, `index :I` | `:*` | Get element at index (0-based). Returns error if out of bounds. |
| `list/set` | `list :L`, `index :I`, `value :*` | `:L` | Set element at index. Returns modified list. Returns error if out of bounds. |
| `list/first` | `list :L` | `:*` | Get first element. Returns error if list is empty. |
| `list/last` | `list :L` | `:*` | Get last element. Returns error if list is empty. |

### Modification

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `list/push` | `list :L`, `element :*` | `:L` | Append element to end of list. Returns modified list. |
| `list/pop` | `list :L` | `:*` | Remove and return last element. Returns error if list is empty. |
| `list/fill` | `list :L`, `value :*` | `:L` | Fill all positions with deep copies of value. Returns modified list. |
| `list/reverse` | `list :L` | `:L` | Reverse list in-place. Returns reversed list. |

### Search

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `list/contains` | `list :L`, `element :*` | `:I` | Check if list contains element. Returns `1` if found, `0` if not. Uses `Encode()` for comparison. |
| `list/index` | `list :L`, `element :*` | `:I` | Find index of first occurrence. Returns `-1` if not found. Uses `Encode()` for comparison. |

### Slicing & Combining

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `list/subset` | `list :L`, `start :I`, `end :I` | `:L` | Copy subset of list (0-indexed, **inclusive** range). Creates deep copies. Returns error if indices out of bounds. |
| `list/slice` | `list :L`, `start :I`, `end :I` | `:L` | Copy slice of list (0-indexed, **exclusive** end). Bounds-safe (auto-clamps). Creates deep copies. |
| `list/concat` | `lists :L...` | `:L` | Concatenate multiple lists into new list (variadic). Creates deep copies of all elements. |
| `list/join` | `list :L`, `separator :S` | `:S` | Join list elements into string with separator. Converts elements using `Encode()` for non-strings. |

### Iteration & Functional

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `list/iter` | `list :L`, `callback :F` | `:I` | Iterate over list, calling callback for each element. Returns `1` if fully iterated, `0` if stopped early. Callback: `(element :*)` → `:I` (1 to continue, 0 to stop). |
| `list/map` | `list :L`, `mapper :F` | `:L` | Create new list by applying function to each element. Mapper: `(element :*)` → `:*`. |
| `list/filter` | `list :L`, `predicate :F` | `:L` | Create new list with elements that pass predicate. Predicate: `(element :*)` → `:I` (1 to include, 0 to exclude). |
| `list/reduce` | `list :L`, `initial :*`, `reducer :F` | `:*` | Reduce list to single value. Reducer: `(accumulator :*`, `element :*)` → `:*`. |

## Type Legend

- `:L` - List
- `:I` - Integer
- `:S` - String
- `:*` - Any type
- `:F` - Function
- `...` - Variadic (accepts multiple arguments)

## Notes

### List Mutation Behavior

**In-Place Modification:**
Functions that modify lists work on the underlying list structure:
- `list/push` - Appends element
- `list/set` - Updates element at index
- `list/fill` - Replaces all elements
- `list/reverse` - Reverses elements

**Important:** Due to Go slice semantics, always reassign the result:
```lisp
(set mylist (list/push mylist 42))
```

**New List Creation:**
Functions that create new lists return independent copies:
- `list/new`, `list/subset`, `list/slice`, `list/concat`
- `list/map`, `list/filter`
- Uses `DeepCopy()` to prevent shared references

### Subset vs. Slice

**`list/subset` - Inclusive Range:**
```lisp
(set lst '(0 1 2 3 4))
(list/subset lst 1 3)  ; (1 2 3) - includes both indices
```
- Both `start` and `end` are included
- Returns error if indices are out of bounds
- Requires valid range (start ≤ end)

**`list/slice` - Exclusive End (Go semantics):**
```lisp
(set lst '(0 1 2 3 4))
(list/slice lst 1 3)   ; (1 2) - excludes end index
```
- `start` is included, `end` is excluded
- Auto-clamps out-of-bounds indices
- Safe for all integer inputs

### Lambda Functions

Functions accepting callbacks (`list/iter`, `list/map`, `list/filter`, `list/reduce`) evaluate the callback function but not before calling it with arguments.

**Inline Lambda:**
```lisp
(list/map mylist (fn (x :I) :I (int/mul x 2)))
```

**Variable Lambda:**
```lisp
(set double (fn (x :I) :I (int/mul x 2)))
(list/map mylist double)
```

### Element Comparison

`list/contains` and `list/index` use `obj.Encode()` for equality comparison:
- Works across all types
- String comparison of encoded representations
- Not reference equality

### Error Handling

Functions that can fail return error objects:
- `list/new` - Negative length
- `list/get`, `list/set` - Index out of bounds
- `list/subset` - Invalid indices or range
- `list/pop`, `list/first`, `list/last` - Empty list
- `list/iter`, `list/map`, `list/filter`, `list/reduce` - Callback errors
- `list/concat` - Non-list arguments

### Variadic Functions

**`list/concat`:**
- Accepts zero or more list arguments
- Returns empty list if no arguments
- Validates all arguments are lists at runtime
- Type mismatch returns error with position

## Examples

### Basic Operations

```lisp
(set lst (list/new 5 0))             ; (0 0 0 0 0)
(putln (list/len lst))                ; 5
(set lst (list/set lst 2 42))        ; (0 0 42 0 0)
(putln (list/get lst 2))              ; 42
(putln (list/first lst))              ; 0
(putln (list/last lst))               ; 0
```

### Building Lists

```lisp
(set nums (uq '()))
(set nums (list/push nums 1))
(set nums (list/push nums 2))
(set nums (list/push nums 3))        ; (1 2 3)

(set last (list/pop nums))           ; last = 3, nums = (1 2)
```

### Iteration

```lisp
(set items '(10 20 30 40 50))

(list/iter items (fn (el :*) :I (do
  (putln el)
  1
)))

(set callback (fn (x :*) :I 
  (if (int/gt x 25)
    0
    (do (putln x) 1))))

(list/iter items callback)           ; Prints 10, 20, stops at 30
```

### Search & Test

```lisp
(set fruits '("apple" "banana" "cherry"))

(putln (list/contains fruits "banana"))      ; 1
(putln (list/index fruits "cherry"))         ; 2
(putln (list/index fruits "grape"))          ; -1
(putln (list/empty fruits))                  ; 0
```

### Slicing & Combining

```lisp
(set nums '(0 1 2 3 4 5 6 7 8 9))

(set sub (list/subset nums 2 5))     ; (2 3 4 5) - inclusive
(set slc (list/slice nums 2 5))      ; (2 3 4) - exclusive end

(set a '(1 2))
(set b '(3 4))
(set c '(5 6))
(set combined (list/concat a b c))   ; (1 2 3 4 5 6)

(set words '("hello" "world"))
(set sentence (list/join words " ")) ; "hello world"
```

### Functional Programming

**Map:**
```lisp
(set nums '(1 2 3 4 5))
(set squared (list/map nums (fn (x :I) :I 
  (int/mul x x))))                   ; (1 4 9 16 25)
```

**Filter:**
```lisp
(set nums '(1 2 3 4 5 6 7 8 9 10))
(set evens (list/filter nums (fn (x :I) :I 
  (int/eq (int/mod x 2) 0))))        ; (2 4 6 8 10)
```

**Reduce:**
```lisp
(set nums '(1 2 3 4 5))
(set sum (list/reduce nums 0 (fn (acc :I el :I) :I 
  (int/add acc el))))                ; 15

(set product (list/reduce nums 1 (fn (acc :I el :I) :I 
  (int/mul acc el))))                ; 120
```

### Advanced Usage

**Find Maximum:**
```lisp
(set nums '(42 17 99 23 56))
(set max (list/reduce nums (list/first nums) 
  (fn (acc :I el :I) :I 
    (if (int/gt el acc) el acc))))   ; 99
```

**List Transformation Pipeline:**
```lisp
(set nums '(1 2 3 4 5 6 7 8 9 10))

(set result (list/map 
  (list/filter nums (fn (x :I) :I (int/gt x 5)))
  (fn (x :I) :I (int/mul x 2))))     ; (12 14 16 18 20)
```

**Reverse and Join:**
```lisp
(set words '("world" "hello"))
(set reversed (list/reverse words))  ; ("hello" "world")
(set sentence (list/join reversed " ")) ; "hello world"
```

**Fill and Modify:**
```lisp
(set grid (list/new 10 0))
(set grid (list/fill grid 1))        ; All 1s
(set grid (list/set grid 5 99))      ; One element changed
```

**Custom Iterator with Early Exit:**
```lisp
(set nums '(1 2 3 4 5 6 7 8 9 10))
(set result (list/iter nums (fn (x :I) :I (do
  (if (int/gt x 5)
    0
    (do
      (putln (int/mul x x))
      1))))))                        ; Prints squares of 1-5, returns 0
```

## Performance Notes

- `list/new` - O(n) for allocation and initialization
- `list/len`, `list/empty`, `list/first`, `list/last` - O(1)
- `list/get`, `list/set` - O(1)
- `list/push` - O(1) amortized (may reallocate)
- `list/pop` - O(1) (does not actually remove, just returns element)
- `list/contains`, `list/index` - O(n) linear search
- `list/reverse` - O(n) in-place swap
- `list/concat` - O(n) where n is total elements
- `list/subset`, `list/slice` - O(n) for range size
- `list/iter`, `list/map`, `list/filter`, `list/reduce` - O(n) iterations
- Deep copies add overhead for complex nested structures

## Implementation Details

**Deep Copy Behavior:**
Functions that create new lists or fill elements use `DeepCopy()` to ensure:
- Nested lists are fully copied
- No shared mutable state
- Safe for concurrent use when properly synchronized

**Iteration Safety:**
Functions like `list/iter`, `list/map`, `list/filter` iterate over a snapshot of the list. Modifications during iteration may not be visible to the iterator.

**Pop Operation:**
`list/pop` returns the last element but does not modify the list slice. To actually remove:
```lisp
(set value (list/pop mylist))
(set mylist (list/slice mylist 0 (int/sub (list/len mylist) 1)))
```

