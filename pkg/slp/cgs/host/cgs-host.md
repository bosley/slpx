# CGS Host Functions (`host`)

Host system information and environment access command group for SLPX.

## Function Reference

### Environment Variables

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `host/env/get` | `name :S` | `:S` | Get environment variable value. Returns error if variable doesn't exist. |
| `host/env/set` | `name :S`, `value :S` | `:I` | Set environment variable. Returns `1` on success, error on failure. |

### System Directories

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `host/dir/home` | | `:S` | Get user home directory. Returns error on failure. |
| `host/dir/config` | | `:S` | Get user configuration directory. Returns error on failure. |
| `host/dir/temp` | | `:S` | Get system temporary directory. |
| `host/dir/cache` | | `:S` | Get user cache directory. Returns error on failure. |

### Operating System

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `host/os` | | `:S` | Get operating system name (e.g., "darwin", "linux", "windows"). |

### Hardware - Memory

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `host/hw/mem/total` | | `:I` | Total memory in bytes. Returns error on failure. |
| `host/hw/mem/available` | | `:I` | Available memory in bytes. Returns error on failure. |
| `host/hw/mem/used` | | `:I` | Used memory in bytes. Returns error on failure. |
| `host/hw/mem/percent` | | `:R` | Memory usage percentage (0-100). Returns error on failure. |

### Hardware - Disk

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `host/hw/disk/total` | | `:I` | Total disk space in bytes. Returns error on failure. |
| `host/hw/disk/used` | | `:I` | Used disk space in bytes. Returns error on failure. |
| `host/hw/disk/percent` | | `:R` | Disk usage percentage (0-100). Returns error on failure. |

### Hardware - CPU (Aggregate)

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `host/hw/cpu/percent` | | `:R` | First CPU usage percentage (0-100). Returns error on failure. |
| `host/hw/cpu/count` | | `:I` | Number of CPUs detected. Returns error on failure. |

### Hardware - CPU (Indexed)

| Function | Parameters | Return Type | Description |
|----------|-----------|-------------|-------------|
| `host/hw/cpu/percent/at` | `idx :I` | `:R` | CPU usage percentage at index (0-100). Returns error if index out of range. |
| `host/hw/cpu/model/at` | `idx :I` | `:S` | CPU model name at index. Returns error if index out of range. |
| `host/hw/cpu/mhz/at` | `idx :I` | `:I` | CPU frequency in MHz at index. Returns error if index out of range. |
| `host/hw/cpu/cache/at` | `idx :I` | `:I` | CPU cache size in bytes at index. Returns error if index out of range. |

## Type Legend

- `:S` - String
- `:I` - Integer (64-bit signed)
- `:R` - Real (64-bit floating-point)

## Notes

### Environment Variables

`host/env/get` and `host/env/set` operate on the **process environment**, not the SLPX REPL environment. Changes made with `host/env/set` affect the running process and any child processes spawned, but do not persist after the program exits.

### Directory Functions

All directory functions return absolute paths. On different operating systems, these paths will vary:

**macOS/Linux:**
- Home: `/Users/username` or `/home/username`
- Config: `$HOME/.config` or system-specific locations
- Cache: `$HOME/.cache` or system-specific locations
- Temp: `/tmp` or system-specific locations

**Windows:**
- Home: `C:\Users\username`
- Config: `%APPDATA%`
- Cache: `%LOCALAPPDATA%`
- Temp: `%TEMP%`

### Hardware Information

Hardware queries use the `gopsutil` library and may take time to gather information. Each call fetches fresh data from the system.

**Memory values** are in bytes. To convert to more readable units:
- KB: divide by 1024
- MB: divide by 1048576 (1024²)
- GB: divide by 1073741824 (1024³)

**Disk information** is for the main disk (typically the home directory's disk on Unix-like systems).

**CPU information** returns data for all logical cores. On systems with hyperthreading, each thread appears as a separate CPU. Use `host/hw/cpu/count` to determine how many CPUs are available, then iterate with the `/at` functions.

### Error Handling

Functions return error objects in the following cases:
- `host/env/get` - Environment variable not found
- `host/env/set` - Failed to set variable (permission issues)
- `host/dir/*` - Failed to determine directory (except `temp`, which never fails)
- `host/hw/*` - Failed to query hardware information
- `host/hw/cpu/*/at` - Index out of range (negative or >= CPU count)

## Examples

### Environment Variables

```lisp
(host/env/get "PATH")
(host/env/get "HOME")

(host/env/set "MY_VAR" "my_value")
(host/env/get "MY_VAR")

(try
  (host/env/get "NONEXISTENT_VAR")
  (putln "Variable not found"))
```

### System Directories

```lisp
(host/dir/home)
(host/dir/config)
(host/dir/temp)
(host/dir/cache)

(set home (host/dir/home))
(putln (str/concat "Home: " home))
```

### Operating System

```lisp
(host/os)

(set os (host/os))
(if (str/eq os "darwin")
  (putln "Running on macOS")
  (if (str/eq os "linux")
    (putln "Running on Linux")
    (if (str/eq os "windows")
      (putln "Running on Windows")
      (putln "Unknown OS"))))
```

### Memory Information

```lisp
(host/hw/mem/total)
(host/hw/mem/available)
(host/hw/mem/used)
(host/hw/mem/percent)

(set total (host/hw/mem/total))
(set used (host/hw/mem/used))
(set percent (host/hw/mem/percent))
(putln (str/concat "Memory: " (str/from used) " / " (str/from total)))
(putln (str/concat "Usage: " (str/from percent) "%"))

(set total_gb (real/div (int/real total) 1073741824.0))
(putln (str/concat "Total Memory: " (str/from total_gb) " GB"))
```

### Disk Information

```lisp
(host/hw/disk/total)
(host/hw/disk/used)
(host/hw/disk/percent)

(set disk_total (host/hw/disk/total))
(set disk_used (host/hw/disk/used))
(set disk_percent (host/hw/disk/percent))

(set available (int/sub disk_total disk_used))
(set avail_gb (real/div (int/real available) 1073741824.0))
(putln (str/concat "Available Disk Space: " (str/from avail_gb) " GB"))
```

### CPU Information (Aggregate)

```lisp
(host/hw/cpu/count)
(host/hw/cpu/percent)

(set cpu_count (host/hw/cpu/count))
(putln (str/concat "CPU Count: " (str/from cpu_count)))

(set cpu_usage (host/hw/cpu/percent))
(putln (str/concat "First CPU Usage: " (str/from cpu_usage) "%"))
```

### CPU Information (Indexed)

```lisp
(set cpu_count (host/hw/cpu/count))
(set idx 0)
(loop (int/lt idx cpu_count)
  (do
    (set percent (host/hw/cpu/percent/at idx))
    (set model (host/hw/cpu/model/at idx))
    (set mhz (host/hw/cpu/mhz/at idx))
    (set cache (host/hw/cpu/cache/at idx))
    
    (putln (str/concat "CPU " (str/from idx) ":"))
    (putln (str/concat "  Model: " model))
    (putln (str/concat "  MHz: " (str/from mhz)))
    (putln (str/concat "  Cache: " (str/from cache) " bytes"))
    (putln (str/concat "  Usage: " (str/from percent) "%"))
    
    (set idx (int/add idx 1))))
```

### Error Handling

```lisp
(try
  (host/hw/cpu/percent/at 999)
  (putln "CPU index out of range"))

(try
  (host/env/get "DOES_NOT_EXIST")
  (putln "Environment variable not found"))
```

### Complete System Profile

```lisp
(putln "=== System Profile ===")
(putln (str/concat "OS: " (host/os)))
(putln (str/concat "Home: " (host/dir/home)))
(putln "")

(putln "=== Memory ===")
(set mem_total (host/hw/mem/total))
(set mem_used (host/hw/mem/used))
(set mem_percent (host/hw/mem/percent))
(putln (str/concat "Total: " (str/from (real/div (int/real mem_total) 1073741824.0)) " GB"))
(putln (str/concat "Used: " (str/from (real/div (int/real mem_used) 1073741824.0)) " GB"))
(putln (str/concat "Usage: " (str/from mem_percent) "%"))
(putln "")

(putln "=== Disk ===")
(set disk_total (host/hw/disk/total))
(set disk_used (host/hw/disk/used))
(set disk_percent (host/hw/disk/percent))
(putln (str/concat "Total: " (str/from (real/div (int/real disk_total) 1073741824.0)) " GB"))
(putln (str/concat "Used: " (str/from (real/div (int/real disk_used) 1073741824.0)) " GB"))
(putln (str/concat "Usage: " (str/from disk_percent) "%"))
(putln "")

(putln "=== CPU ===")
(set cpu_count (host/hw/cpu/count))
(putln (str/concat "Count: " (str/from cpu_count)))
(if (int/gt cpu_count 0)
  (do
    (set model (host/hw/cpu/model/at 0))
    (putln (str/concat "Model: " model))
    (set mhz (host/hw/cpu/mhz/at 0))
    (putln (str/concat "MHz: " (str/from mhz))))
  (none))
```

## Implementation Details

### Hardware Queries

Each hardware query calls `GetHardwareProfile()` which:
- Queries virtual memory statistics
- Queries disk usage for the home directory
- Queries CPU usage with a 1-second sampling period
- Queries detailed CPU information (model, cache, frequency)

This means hardware queries are **not cached** and reflect real-time system state on each call.

### Environment Variable Behavior

- `os.Getenv` / `os.LookupEnv` - Read from process environment
- `os.Setenv` - Write to process environment
- Changes affect the current process and child processes
- Changes do not persist after program exit
- On Unix-like systems, environment variables are case-sensitive
- On Windows, environment variables are case-insensitive

### CPU Indexing

CPU indices are 0-based. Valid indices range from `0` to `(host/hw/cpu/count) - 1`. Attempting to access an invalid index returns an error object.

### Byte Size Conversions

All memory and disk values are in bytes (uint64). For display purposes:

```lisp
(fn bytes-to-kb (bytes) (real/div (int/real bytes) 1024.0))
(fn bytes-to-mb (bytes) (real/div (int/real bytes) 1048576.0))
(fn bytes-to-gb (bytes) (real/div (int/real bytes) 1073741824.0))

(set mem (host/hw/mem/total))
(putln (str/concat (str/from (bytes-to-gb mem)) " GB"))
```

