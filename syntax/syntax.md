# SLPX Syntax Extension Installation Guide

This guide explains how to install the SLPX syntax highlighting extension for Visual Studio Code and Cursor on Windows and macOS.

## Overview

The SLPX syntax extension provides:
- Syntax highlighting for all 115+ SLPX commands
- Color-coded type annotations (`:I`, `:R`, `:S`, `:L`, `:F`, `:E`, `:*`, etc.)
- Comment highlighting (`;` line comments)
- String and number literal highlighting
- Special variable highlighting (`$error`, `$args`, `_`)
- Keyword highlighting (fn, set, if, do, try, match, etc.)

## Installation Methods

### Method 1: Direct Installation (Recommended)

This method works for both VSCode and Cursor on all platforms.

#### Windows

1. Open PowerShell or Command Prompt
2. Navigate to the syntax directory:
   ```powershell
   cd path\to\slpx\syntax
   ```
3. Copy the extension to your extensions directory:
   
   For VSCode:
   ```powershell
   xcopy /E /I . "%USERPROFILE%\.vscode\extensions\slpx-1.0.0"
   ```
   
   For Cursor:
   ```powershell
   xcopy /E /I . "%USERPROFILE%\.cursor\extensions\slpx-1.0.0"
   ```

#### macOS / Linux

1. Open Terminal
2. Navigate to the syntax directory:
   ```bash
   cd /path/to/slpx/syntax
   ```
3. Copy the extension to your extensions directory:
   
   For VSCode:
   ```bash
   mkdir -p ~/.vscode/extensions/slpx-1.0.0
   cp -r * ~/.vscode/extensions/slpx-1.0.0/
   ```
   
   For Cursor:
   ```bash
   mkdir -p ~/.cursor/extensions/slpx-1.0.0
   cp -r * ~/.cursor/extensions/slpx-1.0.0/
   ```

### Method 2: Symbolic Link (Development)

This method creates a symbolic link, allowing you to edit the syntax files and see changes immediately after reloading the window.

#### Windows (Administrator PowerShell required)

For VSCode:
```powershell
New-Item -ItemType SymbolicLink -Path "$env:USERPROFILE\.vscode\extensions\slpx-1.0.0" -Target "path\to\slpx\syntax"
```

For Cursor:
```powershell
New-Item -ItemType SymbolicLink -Path "$env:USERPROFILE\.cursor\extensions\slpx-1.0.0" -Target "path\to\slpx\syntax"
```

#### macOS / Linux

For VSCode:
```bash
ln -s /path/to/slpx/syntax ~/.vscode/extensions/slpx-1.0.0
```

For Cursor:
```bash
ln -s /path/to/slpx/syntax ~/.cursor/extensions/slpx-1.0.0
```

### Method 3: VSIX Package (Distribution)

To create a distributable package:

1. Install vsce (VSCode Extension Manager):
   ```bash
   npm install -g @vscode/vsce
   ```

2. Navigate to the syntax directory and package:
   ```bash
   cd /path/to/slpx/syntax
   vsce package
   ```

3. Install the generated `.vsix` file:
   
   In VSCode/Cursor:
   - Press `Cmd+Shift+P` (macOS) or `Ctrl+Shift+P` (Windows)
   - Type "Extensions: Install from VSIX"
   - Select the generated `slpx-1.0.0.vsix` file

## Verifying Installation

1. Restart VSCode or Cursor (or reload the window: `Cmd+R` on macOS, `Ctrl+R` on Windows)
2. Open any `.slpx` file
3. Check the bottom-right corner of the editor - it should show "SLPX" as the language
4. Verify syntax highlighting is active:
   - Keywords like `fn`, `set`, `if` should be highlighted
   - Function names like `int/add`, `list/map` should be highlighted
   - Comments starting with `;` should be styled as comments
   - Strings in `"quotes"` should be highlighted

## Manual Language Selection

If a `.slpx` file doesn't automatically use the SLPX syntax:

1. Click the language indicator in the bottom-right corner
2. Type "SLPX" in the search box
3. Select "SLPX" from the list

## Troubleshooting

### Extension Not Loading

1. Check that the extension directory exists:
   - VSCode: `~/.vscode/extensions/slpx-1.0.0/` (macOS/Linux) or `%USERPROFILE%\.vscode\extensions\slpx-1.0.0\` (Windows)
   - Cursor: `~/.cursor/extensions/slpx-1.0.0/` (macOS/Linux) or `%USERPROFILE%\.cursor\extensions\slpx-1.0.0\` (Windows)

2. Verify the extension structure:
   ```
   slpx-1.0.0/
   ├── package.json
   ├── language-configuration.json
   └── syntaxes/
       └── slpx.tmLanguage.json
   ```

3. Check the Developer Tools console for errors:
   - Press `Cmd+Option+I` (macOS) or `Ctrl+Shift+I` (Windows)
   - Look for any errors related to "slpx"

### Syntax Not Highlighting

1. Reload the window: `Cmd+R` (macOS) or `Ctrl+R` (Windows)
2. Manually select the SLPX language (see "Manual Language Selection" above)
3. Check that the file has a `.slpx` extension

### Colors Look Wrong

Your color theme may not support all token types. Try switching to a different theme:
1. Press `Cmd+K Cmd+T` (macOS) or `Ctrl+K Ctrl+T` (Windows)
2. Select a different theme (e.g., "Dark+ (default dark)")

## Customizing Colors

To customize syntax highlighting colors, add to your VSCode/Cursor `settings.json`:

```json
{
  "editor.tokenColorCustomizations": {
    "textMateRules": [
      {
        "scope": "keyword.control.slpx",
        "settings": {
          "foreground": "#C586C0",
          "fontStyle": "bold"
        }
      },
      {
        "scope": "support.function.slpx",
        "settings": {
          "foreground": "#DCDCAA"
        }
      },
      {
        "scope": "storage.type.slpx",
        "settings": {
          "foreground": "#4EC9B0"
        }
      }
    ]
  }
}
```

## Uninstalling

To remove the extension:

### Windows

For VSCode:
```powershell
Remove-Item -Recurse -Force "$env:USERPROFILE\.vscode\extensions\slpx-1.0.0"
```

For Cursor:
```powershell
Remove-Item -Recurse -Force "$env:USERPROFILE\.cursor\extensions\slpx-1.0.0"
```

### macOS / Linux

For VSCode:
```bash
rm -rf ~/.vscode/extensions/slpx-1.0.0
```

For Cursor:
```bash
rm -rf ~/.cursor/extensions/slpx-1.0.0
```

Then reload the window.

## Additional Resources

- SLPX Command Reference: See `commands.md` for complete list of all 115 commands
- SLPX Type System: `:I` (integer), `:R` (real), `:S` (string), `:L` (list), `:F` (function), `:E` (error), `:*` (any), `:_` (none)
- Issues or improvements: Contact bosley at insula.dev

