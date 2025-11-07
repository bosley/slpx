# SLPX Syntax Extension for VSCode and Cursor

Official syntax highlighting extension for the SLPX programming language.

## What's Included

This extension provides comprehensive syntax highlighting for SLPX, including:

- **115+ Built-in Commands** across 7 command groups:
  - `bits` - Bit manipulation (3 commands)
  - `fs` - File system operations (14 commands)
  - `io` - Input/output and colors (9 commands)
  - `list` - List operations and functional programming (23 commands)
  - `numbers` - Integer and real arithmetic (38 commands)
  - `reflection` - Type introspection (11 commands)
  - `str` - String manipulation (17 commands)

- **Core Language Keywords**: `fn`, `set`, `if`, `do`, `try`, `match`, `use`, `exit`, `drop`, `qu`, `uq`, `putln`

- **Type Annotations**: `:I`, `:R`, `:S`, `:L`, `:F`, `:E`, `:*`, `:_`, `:Q`, `:X`

- **Special Variables**: `$error`, `$args`, `_`

- **Syntax Elements**: Comments (`;`), strings, numbers, parentheses, quotes

## Directory Structure

```
syntax/
├── README.md                      # This file
├── syntax.md                      # Installation instructions
├── commands.md                    # Complete command reference
├── package.json                   # Extension manifest
├── language-configuration.json    # Language features (brackets, comments)
└── syntaxes/
    └── slpx.tmLanguage.json      # TextMate grammar definition
```

## Quick Install

See `syntax.md` for detailed installation instructions for Windows, macOS, VSCode, and Cursor.

### Quick Command

**macOS/Linux (VSCode):**
```bash
mkdir -p ~/.vscode/extensions/slpx-1.0.0 && cp -r * ~/.vscode/extensions/slpx-1.0.0/
```

**macOS/Linux (Cursor):**
```bash
mkdir -p ~/.cursor/extensions/slpx-1.0.0 && cp -r * ~/.cursor/extensions/slpx-1.0.0/
```

Then reload your editor window.

## Features

### Intelligent Highlighting

The extension recognizes and highlights:

1. **Keywords and Control Flow**
   ```slpx
   (fn (x :I) :I (int/add x 1))
   (set result (if (int/gt x 0) x (int/sub 0 x)))
   ```

2. **All Command Groups**
   ```slpx
   (list/map numbers (fn (n :I) :I (int/mul n 2)))
   (fs/write_file "output.txt" (str/concat "Result: " result))
   ```

3. **Type Annotations**
   ```slpx
   (fn (name :S count :I) :L ...)
   ```

4. **Comments and Documentation**
   ```slpx
   ; This is a comment
   ; Comments are styled distinctly
   ```

### Auto-Completion Support

The extension includes:
- Auto-closing parentheses and brackets
- Auto-closing quotes
- Bracket matching
- Comment toggling (`Cmd+/` or `Ctrl+/`)

## Customization

You can customize the colors in your VSCode/Cursor settings:

1. Open Settings (`Cmd+,` or `Ctrl+,`)
2. Search for "token color customizations"
3. Edit the `editor.tokenColorCustomizations` setting

Example:
```json
{
  "editor.tokenColorCustomizations": {
    "textMateRules": [
      {
        "scope": "keyword.control.slpx",
        "settings": {
          "foreground": "#C586C0"
        }
      }
    ]
  }
}
```

## Grammar Scopes

The extension defines these TextMate scopes:

- `comment.line.semicolon.slpx` - Comments
- `string.quoted.double.slpx` - String literals
- `constant.numeric.integer.slpx` - Integer numbers
- `constant.numeric.float.slpx` - Real numbers
- `storage.type.slpx` - Type annotations
- `keyword.control.slpx` - Core keywords
- `support.function.*.slpx` - Command group functions
- `variable.language.slpx` - Special variables
- `constant.language.slpx` - Constants

## Compatibility

- **VSCode**: Version 1.60.0 or higher
- **Cursor**: All versions
- **Platforms**: Windows, macOS, Linux

## Development

To modify the syntax highlighting:

1. Edit `syntaxes/slpx.tmLanguage.json`
2. Reload your editor window (`Cmd+R` or `Ctrl+R`)
3. Test with `.slpx` files

For symlinked installations, changes are reflected immediately after reload.

## Language Server Protocol (Future)

This extension currently provides syntax highlighting only. Future enhancements may include:
- IntelliSense/autocomplete
- Go to definition
- Hover documentation
- Diagnostics/linting
- Code formatting

## Support

For issues or feature requests related to the SLPX language or this extension, contact bosley at insula.dev.

## Version History

### 1.0.0
- Initial release
- Complete syntax highlighting for all 115+ commands
- Support for all type annotations
- Comment and string highlighting
- Auto-closing pairs and bracket matching

