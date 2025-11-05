# CGS Filesystem Functions (`fs`)

Filesystem operations command group for SLPX. Provides comprehensive file and directory manipulation capabilities with proper error handling.

## Function Reference

### Path Predicates

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `fs/exists?` | `path :S` | `:I` | Returns `1` if path exists, `0` otherwise. |
| `fs/dir?` | `path :S` | `:I` | Returns `1` if path is a directory, `0` otherwise. |
| `fs/file?` | `path :S` | `:I` | Returns `1` if path is a file, `0` otherwise. |

### File Operations

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `fs/read_file` | `path :S` | `:S` | Reads and returns file contents as a string. Returns error on failure. |
| `fs/write_file` | `path :S`, `data :S` | `:I` | Writes data to file (overwrites existing). Returns `1` on success, error on failure. |
| `fs/append_file` | `path :S`, `data :S` | `:I` | Appends data to file (creates if doesn't exist). Returns `1` on success, error on failure. |
| `fs/rm_file` | `path :S` | `:I` | Removes a file. Returns `1` on success, error on failure. |

### Directory Operations

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `fs/mk_dir` | `path :S` | `:I` | Creates a single directory. Returns `1` on success, error on failure. |
| `fs/mk_dir_all` | `path :S` | `:I` | Creates directory and all parent directories. Returns `1` on success, error on failure. |
| `fs/rm_dir` | `path :S` | `:I` | Removes an empty directory. Returns `1` on success, error on failure. |
| `fs/rm_dir_all` | `path :S` | `:I` | Removes directory and all contents recursively. Returns `1` on success, error on failure. |
| `fs/list_dir` | `path :S` | `:L` | Returns list of filenames in directory as strings. Returns error on failure. |

### Working Directory

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `fs/working_dir` | none | `:S` | Returns current working directory path. |
| `fs/set_working_dir` | `path :S` | `:I` | Changes working directory. Returns `1` on success, error on failure. |

## Type Legend

- `:S` - String
- `:I` - Integer (64-bit signed)
- `:L` - List
- `:E` - Error

## Notes

### Path Resolution

All paths are resolved relative to the current working directory unless absolute:
- Relative paths: `"data.txt"`, `"subdir/file.txt"`
- Absolute paths: `"/tmp/test.txt"`, `"/home/user/docs"`

The working directory is initially set to the directory containing the executed script.

### Return Values

Functions return:
- `1` (integer) for successful operations
- Error object for failures

Predicate functions (`exists?`, `dir?`, `file?`) return:
- `1` for true
- `0` for false

### Error Handling

All filesystem operations that can fail return error objects with descriptive messages:
```lisp
(fs/read_file "nonexistent.txt")
; Returns: @(failed to read file: no such file or directory)
```

Use error handling constructs to catch failures:
```lisp
(if (reflect/error? (fs/read_file "config.txt"))
  (putln "Failed to read config")
  (putln "Config loaded"))
```

### File Permissions

- Files created with `fs/write_file` and `fs/append_file` have permissions `0644` (rw-r--r--)
- Directories created with `fs/mk_dir` and `fs/mk_dir_all` have permissions `0755` (rwxr-xr-x)

## Examples

### Path Checking

```lisp
(fs/exists? "/tmp/data.txt")           ; 1 if exists, 0 otherwise
(fs/dir? "/tmp")                       ; 1 (directory)
(fs/file? "/tmp/data.txt")             ; 1 (file)
(fs/file? "/tmp")                      ; 0 (not a file)
```

### Reading Files

```lisp
(set content (fs/read_file "config.txt"))
(if (reflect/error? content)
  (putln "Failed to read file")
  (putln content))
```

### Writing Files

```lisp
(fs/write_file "output.txt" "Hello, World!")
; Returns 1 on success

(fs/write_file "output.txt" "New content")
; Overwrites existing content
```

### Appending to Files

```lisp
(fs/write_file "log.txt" "Line 1\n")
(fs/append_file "log.txt" "Line 2\n")
(fs/append_file "log.txt" "Line 3\n")
; log.txt now contains:
; Line 1
; Line 2
; Line 3
```

### Creating Directories

```lisp
(fs/mk_dir "data")
; Creates single directory

(fs/mk_dir_all "project/data/cache")
; Creates entire directory tree
```

### Removing Files and Directories

```lisp
(fs/rm_file "temp.txt")
; Removes single file

(fs/rm_dir "empty_dir")
; Removes empty directory only

(fs/rm_dir_all "project")
; Removes directory and all contents
```

### Listing Directory Contents

```lisp
(set files (fs/list_dir "/tmp"))
; Returns list: ("file1.txt" "file2.txt" "subdir" ...)

(list/each files (fn (file) :_
  (putln file)))
; Prints each filename
```

### Working Directory Management

```lisp
(putln (fs/working_dir))
; Prints current directory

(fs/set_working_dir "/tmp")
; Change to /tmp

(fs/write_file "test.txt" "data")
; Creates /tmp/test.txt
```

### Safe File Operations

```lisp
(fn (safe_read_file path) :S
  (if (fs/exists? path)
    (if (fs/file? path)
      (fs/read_file path)
      '@(path is not a file))
    '@(file does not exist)))

(set result (safe_read_file "config.txt"))
(if (reflect/error? result)
  (putln "Error reading file")
  (putln result))
```

### Directory Tree Creation

```lisp
(fn (ensure_dir path) :I
  (if (fs/exists? path)
    (if (fs/dir? path)
      1
      '@(path exists but is not a directory))
    (fs/mk_dir_all path)))

(ensure_dir "output/data/cache")
```

### Conditional File Writing

```lisp
(fn (write_if_not_exists path data) :*
  (if (fs/exists? path)
    '@(file already exists)
    (fs/write_file path data)))

(write_if_not_exists "output.txt" "Initial content")
```

### Directory Listing with Filtering

```lisp
(set files (fs/list_dir "."))
(set txt_files (list/filter files (fn (name) :I
  (str/ends_with? name ".txt"))))

(list/each txt_files (fn (file) :_
  (putln (str/concat "Found: " file))))
```

### Temporary Directory Pattern

```lisp
(set temp_dir "/tmp/my-app-temp")

(do
  (if (fs/exists? temp_dir)
    (fs/rm_dir_all temp_dir)
    _)
  
  (fs/mk_dir temp_dir)
  
  (fs/write_file 
    (str/concat temp_dir "/data.txt")
    "temporary data")
  
  (fs/rm_dir_all temp_dir))
```

### Copy File Pattern

```lisp
(fn (copy_file src dest) :*
  (if (fs/file? src)
    (do
      (set content (fs/read_file src))
      (if (reflect/error? content)
        content
        (fs/write_file dest content)))
    '@(source is not a file)))

(copy_file "input.txt" "backup.txt")
```

### Recursive Directory Walk Pattern

```lisp
(fn (process_directory path) :_
  (set entries (fs/list_dir path))
  (if (reflect/error? entries)
    (putln (str/concat "Error: " path))
    (list/each entries (fn (name) :_
      (set full_path (str/concat path "/" name))
      (if (fs/dir? full_path)
        (process_directory full_path)
        (putln (str/concat "File: " full_path)))))))

(process_directory ".")
```

## Performance Notes

- Path resolution is efficient using `filepath.Join`
- All operations use OS-level filesystem calls
- `fs/list_dir` returns only filenames, not full metadata
- `fs/rm_dir_all` is recursive and can be slow for large directory trees
- Relative paths are always resolved against working directory

## Implementation Details

**Path Resolution:**
- All paths resolved through internal `resolvePath` method
- Absolute paths used as-is
- Relative paths joined with current working directory

**Working Directory:**
- Initially set from script's directory via `Runtime.GetStartPath()`
- Changed via `fs/set_working_dir` affects all subsequent operations
- Working directory state persists across function calls

**Error Handling:**
- All errors return `object.OBJ_TYPE_ERROR` with descriptive messages
- File operations preserve original OS error details
- Predicate functions never error, always return `0` or `1`

**File Permissions:**
- Write operations: `0644` (owner: rw, group: r, other: r)
- Directory operations: `0755` (owner: rwx, group: rx, other: rx)

