# README

This readme covers a projec I've been considering making in go for quite some time. Its just a simple list processor.

This document fully details how it works, but if you're impatient and just want to see the code, run:

```bash
make
./build/slpx
```

You will be prompted to setup the environment (a double-press of return should do it to accept defaults)

Then, you can run:

```bash
 ./build/slpx tests/primitive/main.slpx 
```

To see the main test suite run, or (recommended) just run:

```bash
 ./build/slpx
```

This will launcht the TUI REPL environment. This environment is a little basic, but much more than your average
language demo REPL. 

## Examples/Etc

In `examples/` you will find runnable samples that you can run as a file into the main `slpx` binary, or you can
pick and mix into the REPL editor to explore the environment.

For more advanced examples and a very wide breadth of all commands used, see `tests/primitive/bootstrap.slpx` and
examine its role in the `primitive` tests by reading the `tests/primitives/main.slpx` in relation tot he bootstrap
file.

# SLP - Parser & Data

Parses test into lists of the following:

- integer
- real
- string
- error
- lists
- identifier
- none
- some (aka quoted)

## Integer

All integer numbers, signed or unsigned in base 10. 
Represented in Go by a 64-bit integer that preserves the encoded sign.

## Real

Real numbers are detected via the presence of a single `.` following and/or
preceding the presence of an integer.
Represented in Go by a 64-bit floating point.

## String

Whenever the parser detects a `"` ASCII character, it goes into a "parsing string" context that won't complete until it finds
a matching `"` that isn't directly preceded by a `\` (an escape char). While operating in the parsing string context, the
escape char serves as an instruction to "not leave the string" and to "accept" the `"` following the `\` to be part of the
string being read in. This means that the `\` prior to `"` in a `\"` occurrence will be dropped and not manifest
itself in the represented data at runtime directly.

## Identifier

Identifiers are any unmatched grouping of data between sets of ` ` whitespace within some context (like a `()` [see below]). 
Labeling these as `identifier` presupposes that there will be something doing the identifying. There was some consideration
in labeling it `collection` or `raw`, but I figured that the existence of a parser itself implies the existence of something
that wants to do parsing. Sticking with the historically-used `identifier` implies that someone somewhere will want to 
identify it and will somehow have context to make sense of it at some point (the runtime/env).

## Error

Not an error from parsing, but a physical manifestation indicating that something, somewhere, went wrong. The SLP
identifies any list prefixed with a `@` symbol and parses it into an `error object`. Parsing one of these is perfectly fine,
and much like calling a group `identifier`, we assume someone somewhere (the runtime/env) will know how to handle it in the context
in which they observe it.

## List

A list in SLP is defined as "a collection of parsed objects" that are inscribed using a pair of parentheses `()`.
A list can contain any of the listed objects, even other lists (of course).

## Some

Aka a "quoted" is "any valid parsed object that follows a `'` symbol." This is useful for the environment by permitting a
means to "delay evaluation" and a _lot_ more. 

```
Note on verbiage :

A "quoted" item is seen in the form `'a`, `'3.14`, `'(1 2 3)` as it is the visual, readable representation of concept, 
but sometimes the word `quoted` is used to describe this concept in the `context of runtime evaluation`, when it's technically
a `some` object.
```

## None

A "none" is "nothing." It can be considered "the proper lack of an object." This is whenever `_` is detected during the
parsing stage (a single-length identifier of identification value `_`, essentially.)

---

## Macros

Macros are a parse-time template expansion mechanism that enable code generation and syntactic abstraction. Unlike runtime identifiers beginning with `$`, macros operate during the parsing phase and use the `$` symbol for both definition and invocation.

### Definition

A macro is defined using the syntax `$(pattern) template`, where `pattern` is a list containing the macro name and zero or more parameters, and `template` is any valid SLP expression.

```slpx
$(macro_name ?param1 ?param2) template_expression
```

The macro name must be an identifier. Parameters are identifiers prefixed with `?`. The template can be any valid parsed object including integers, strings, lists, quoted expressions, or nested structures. After a macro definition is parsed, it returns `none` and is stored in the parser's macro table for subsequent expansions.

### Invocation

A macro is invoked by prefixing its name with `$` as the first element of a list.

```slpx
($macro_name arg1 arg2)
```

When the parser encounters a list beginning with an identifier prefixed by `$`, it checks the macro table for a matching definition. If found, the macro's template is expanded by substituting each occurrence of a parameter with the corresponding argument. The number of arguments must match the number of parameters exactly or a parse error will result.

### Template Substitution

During expansion, each parameter identifier (those prefixed with `?` in the pattern) is replaced with the corresponding argument provided at the call site. This substitution occurs recursively throughout the template structure.

```slpx
$(identity ?x) ?x
($identity 42)
```

The above expands to `42` at parse time. Parameters can appear multiple times within a template, and each occurrence is substituted with the argument.

```slpx
$(twice ?x) (if ?x ?x 0)
($twice 5)
```

This expands to `(if 5 5 0)`. Substitution operates on all nested structures including lists and quoted expressions.

```slpx
$(wrap ?val) (qu (?val nested ?val))
($wrap test)
```

This expands to `(qu (test nested test))`. Arguments are deep-copied during substitution to prevent unintended mutations.

### Expansion Timing

Macro expansion occurs during parsing via the `expandMacroIfNeeded` function, which is called after a list is successfully parsed. This means macros are expanded before runtime evaluation begins, enabling them to generate arbitrary code structures that are then evaluated normally.

Macro expansion is recursive. If a macro expands to a list that itself begins with a macro call, the expansion process continues until no further macro calls are detected.

```slpx
$(inner ?v) (if ?v ?v 0)
$(outer ?v) ($inner ?v)
($outer 25)
```

The `outer` macro expands to `($inner 25)`, which then expands to `(if 25 25 0)`.

### Scope and Redefinition

Macros are scoped to the parser instance. A macro defined in one file is available in subsequent parsing within that same session, but is not automatically available across separate parse invocations unless explicitly re-defined or imported.

Macro definitions can be redefined. A subsequent definition with the same name replaces the previous definition in the macro table.

```slpx
$(test ?x) ?x
$(test ?y) (if ?y ?y 0)
($test 1)
```

The second definition of `test` overwrites the first, and the invocation uses the most recent definition.

### Error Conditions

The parser enforces several constraints on macro definitions and invocations:

**Undefined macro**: Attempting to invoke a macro that has not been defined results in a parse error indicating the macro name is undefined.

**Arity mismatch**: The number of arguments provided at the call site must exactly match the number of parameters in the macro's pattern. Providing too few or too many arguments results in a parse error.

**Empty pattern**: A macro definition must have at least a name. An empty pattern list `$() template` results in a parse error.

**Non-identifier in pattern**: All elements of the pattern list must be identifiers. If a non-identifier (such as an integer or string) appears in the pattern, a parse error occurs.

**Parameter without `?` prefix**: All parameters after the macro name must be identifiers beginning with `?`. If an identifier in the parameter position does not start with `?`, a parse error occurs.

### Use Cases

Macros enable several patterns:

**Code generation**: Create functions or other structures programmatically.

```slpx
$(defun ?name ?body) (set ?name (fn () :I ?body))
($defun my_func 123)
```

This generates `(set my_func (fn () :I 123))` at parse time.

**Conditional abstractions**: Build higher-level control flow constructs.

```slpx
$(when ?cond ?body) (if ?cond ?body _)
($when 1 (putln "executed"))
```

Expands to `(if 1 (putln "executed") _)`.

**Repeated patterns**: Eliminate boilerplate by capturing common structures.

```slpx
$(repeat3 ?expr) (do ?expr ?expr ?expr)
($repeat3 (putln "hello"))
```

Expands to `(do (putln "hello") (putln "hello") (putln "hello"))`.

**DSL construction**: Define domain-specific syntax that expands to core language constructs, enabling more expressive or specialized notation for particular problem domains.

---

# Runtime

The runtime is an environment (see pkg/slp/env) of an `IO` `FS` and `MEM` interface along with some logic to
express what it means to "process" an SLP list. These 3 items are defined and handed to our runtime so we can section-off
and re-define the implementation behind how to load and store data mapped to identifiers.

This means we can leverage a virutal file system, and set hard and very controllable upper limits on activity for any
given script/repl enviornment.

## Commands

A command is a function implemented by the runtime that can be triggered by pre-set identifiers during evaluation time. 
The SLP runtime operates on two categories of commands:

### Core Commands

These commands are found in the slp `env` directly, as the core building blocks of the language. These commands contain
idendtifiers `set` `drop` `fn`, and more. 

### Command Grouped Symbols

`CGS` are groups of symbols defined in such a way that they can be "injected" into a runtime "in addition to" the "core"
commands. Most of the existing commands are part of a CGS. See pkg/slp/cgs to see what commands are central enough to
be grouped in the language itself, but not technically the `most central` set of commands.

## Type Symbols

These symbols are used by the runtime when parsing function definitions (the `fn` command) and/or a matching
call-site, to streamline list passing expectations when invoking functions. There is not yet a typecheck pre-flight.

| Symbol | Type       |
|--------|------------|
| :_     | none       |
| :Q     | some       |
| :*     | any        |
| :L     | list       |
| :E     | error      |
| :S     | string     |
| :I     | integer    |
| :R     | real       |
| :X     | identifier |
| :F     | function   |

## Variadics

Use `..` as the parameter list to create variadic functions that accept any number of arguments. All arguments are evaluated and available as `$args` inside the function body.

```slpx
(set sum (fn (..) (do-something-with $args)))
(sum 1 2 3 4 5)
```

Built-in variadic functions like `putln`, `use`, `do`, `int/sum`, and `real/sum` require at least one argument.

When a user-defined variadic function is invoked, the runtime evaluates each argument and constructs a list object that's injected into the function's local memory scope as `$args`. This injection happens automatically before the function body executes and is scoped to the function's execution context, meaning `$args` is not available outside of the function call.

## System-Reserved Identifiers

Identifiers prefixed with `$` are reserved exclusively for runtime use and cannot be defined by user code. This restriction is enforced at the time of assignment via the `set` command, which will return an error if an attempt is made to define an identifier beginning with `$`. 

The runtime leverages this reserved namespace to inject context-specific identifiers into evaluation scopes without risk of collision with user-defined symbols. Currently, two system identifiers exist:

**`$args`** is injected into the local memory scope of variadic functions. It contains a list of all evaluated arguments passed to the function. Once the function completes execution, `$args` is no longer accessible.

**`$error`** is injected into the handler body of a `try` statement when the attempted expression results in an error. The `$error` identifier contains the error message as a string. After the handler completes, `$error` is explicitly removed from memory and is no longer available.

This design permits the runtime to provide contextual data to executing code while maintaining a clear separation between user space and system space.
