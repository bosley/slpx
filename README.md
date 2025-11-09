# SLPX - Simple List Processor

A configurable list-processing language with a TUI REPL environment, macro system, and extensible runtime.

## Table of Contents

- [Quick Start](#quick-start)
- [TUI Interface](#tui-interface)
- [Examples/Etc](#examplesetc)
- [Customization](#customization)
- [Syntax Highlighting](#syntax-highlighting)
- [SLP - Parser & Data](#slp---parser--data)
  - [Integer](#integer)
  - [Real](#real)
  - [String](#string)
  - [Identifier](#identifier)
  - [Error](#error)
  - [List](#list)
  - [Some](#some)
  - [None](#none)
  - [Macros](#macros)
    - [Definition](#definition)
    - [Invocation](#invocation)
    - [Template Substitution](#template-substitution)
    - [Expansion Timing](#expansion-timing)
    - [Scope and Redefinition](#scope-and-redefinition)
    - [Error Conditions](#error-conditions)
    - [Use Cases](#use-cases)
- [Runtime](#runtime)
  - [Commands](#commands)
    - [Core Commands](#core-commands)
    - [Command Grouped Symbols](#command-grouped-symbols)
      - [Available Command Groups](#available-command-groups)
        - [Bits](pkg/slp/cgs/bits/cfgs-bits.md)
        - [Filesystem](pkg/slp/cgs/fs/cgs-fs.md)
        - [Host](pkg/slp/cgs/host/cgs-host.md)
        - [IO](pkg/slp/cgs/io/cgs-io.md)
        - [List](pkg/slp/cgs/list/cgs-list.md)
        - [Numbers](pkg/slp/cgs/numbers/cgs-numbers.md)
        - [Reflection](pkg/slp/cgs/reflection/cgs-reflection.md)
        - [String](pkg/slp/cgs/str/cgs-str.md)
  - [Type Symbols](#type-symbols)
  - [Variadics](#variadics)
  - [System-Reserved Identifiers](#system-reserved-identifiers)
- [Function Execution Architecture](#function-execution-architecture)
  - [Key Architectural Points](#key-architectural-points)
- [Tests](#tests)
  - [Primary SLPX tests](#primary-slpx-tests)
  - [Primitive Test Process](#primitive-test-process)

---

## Quick Start

This document fully details how the language works, but if you're impatient and just want to see the code, run:

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

## TUI Interface

Upon first run, you'll be greeted with an installation prompt to setup the files required to run
the configurable user repl environment. 

<img src="resources/install.png" width="600" alt="Installation">

The editor provides a clean interface for writing SLPX code, able to be toggled from the more traditional
REPL input screen. This screen also contains a vertival list of past commands able to be selected from
and inserted into the editor (not shown.)

<img src="resources/editor.png" width="600" alt="Editor">

The output screen displays execution results with proper output highlighting:

<img src="resources/tui-out.png" width="600" alt="Output">

## Examples/Etc

In `examples/` you will find runnable samples that you can run as a file into the main `slpx` binary, or you can
pick and mix into the REPL editor to explore the environment.

For more advanced examples and a very wide breadth of all commands used, see `tests/primitive/bootstrap.slpx` and
examine its role in the `primitive` tests by reading the `tests/primitives/main.slpx` in relation tot he bootstrap
file.

## Customization

The tui controls and colors can be customized by modifying your `init.slpx` in the operating system's config dir under
"slpx." This dir is commonly `~/.config/` on linux and `/Users/<username>/Library/Application\ Support/` on mac.

The files themselves are source from `cmd/slpx/assets` under `advanced` and `default`.

## Syntax Highlighting

Read `syntax/README.md` to see how to install the syntax files for `.slpx` extensions in VSCode + derivatives.

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

#### Available Command Groups

Detailed documentation for each command group:

- **[Bits](pkg/slp/cgs/bits/cfgs-bits.md)** - Bit-level manipulation and binary conversion functions
- **[Filesystem](pkg/slp/cgs/fs/cgs-fs.md)** - File and directory operations, path manipulation
- **[Host](pkg/slp/cgs/host/cgs-host.md)** - System information, environment variables, hardware queries
- **[IO](pkg/slp/cgs/io/cgs-io.md)** - Input/output operations, color formatting, console interaction
- **[List](pkg/slp/cgs/list/cgs-list.md)** - List manipulation, iteration, and functional programming
- **[Numbers](pkg/slp/cgs/numbers/cgs-numbers.md)** - Arithmetic operations, comparisons, and math functions
- **[Reflection](pkg/slp/cgs/reflection/cgs-reflection.md)** - Type introspection and runtime type checking
- **[String](pkg/slp/cgs/str/cgs-str.md)** - String manipulation, conversion, and processing

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

---

# Function Execution Architecture

The SLP runtime processes functions through a layered evaluation architecture that distinguishes between user-defined functions and runtime-provided functions. The following diagram illustrates the complete execution flow from source text through evaluation to final execution.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              SOURCE TEXT                                    │
│                          "(putln (fn () :I 42))"                            │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │
                                 v
┌─────────────────────────────────────────────────────────────────────────────┐
│                         PARSER (slp/slp.go)                                 │
│  - Tokenizes source into objects                                            │
│  - Expands macros via expandMacroIfNeeded()                                 │
│  - Returns: List, Integer, Real, String, Identifier, Some, None, Error      │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │
                                 v
┌─────────────────────────────────────────────────────────────────────────────┐
│                     EVALUATION CONTEXT (env/env.go)                         │
│                                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                       │
│  │     MEM      │  │      FS      │  │      IO      │                       │
│  │  (memory)    │  │ (filesystem) │  │   (stdio)    │                       │
│  └──────────────┘  └──────────────┘  └──────────────┘                       │
│                                                                             │
│  ┌────────────────────────────────────────────────────────────────┐         │
│  │              FUNCTION GROUP REGISTRY                           │         │
│  │                                                                │         │
│  │  CORE (env/core.go)                                            │         │
│  │    set, putln, fn, try, do, drop, qu, uq, use, exit, if,       │         │
│  │    match                                                       │         │
│  │                                                                │         │
│  │  CGS (pkg/slp/cgs/*)                                           │         │
│  │    - host:       env/get, os, hw/mem/total, hw/cpu/count...    │         │
│  │    - fs:         exists?, read_file, write_file, list_dir...   │         │
│  │    - bits:       explode, int, real                            │         │
│  │    - str:        string operations                             │         │
│  │    - io:         I/O operations                                │         │
│  │    - list:       list operations                               │         │
│  │    - numbers:    numeric operations                            │         │
│  │    - reflection: type introspection                            │         │
│  └────────────────────────────────────────────────────────────────┘         │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │
                                 v
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Evaluate(obj Obj)                                  │
│                                                                             │
│  Switch on obj.Type:                                                        │
│    - NONE, STRING, INTEGER, REAL, ERROR, FUNCTION  →  return as-is          │
│    - SOME (quoted)  →  return without evaluation                            │
│    - IDENTIFIER     →  lookupIdentifier() [check MEM, then FunctionGroups]  │
│    - LIST           →  Execute(list) ↓                                      │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │
                                 v
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Execute(list List)                                 │
│                                                                             │
│  1. Evaluate first element of list                                          │
│  2. Determine callable type                                                 │
└────────────────┬───────────────────────────────────┬────────────────────────┘
                 │                                   │
        ┌────────v───────────┐               ┌───────v─────────────┐
        │  OBJ_TYPE_FUNCTION │               │ OBJ_TYPE_IDENTIFIER │
        │  (User Function)   │               │  (Env Function)     │
        └────────┬───────────┘               └───────┬─────────────┘
                 │                                   │
                 v                                   v
┌─────────────────────────────────────┐  ┌─────────────────────────────────────┐
│  executeObjectFunction()            │  │  executeEnvFunction()               │
│                                     │  │                                     │
│  Source: (fn) command               │  │  Source: FunctionGroup lookup       │
│  Storage: MEM (user variables)      │  │  Storage: functionGroups map        │
│  Closure: Captured MEM context      │  │  Closure: N/A                       │
│                                     │  │                                     │
│  ┌───────────────────────────────┐  │  │  ┌───────────────────────────────┐  │
│  │ 1. Check if Variadic          │  │  │  │ 1. Evaluate Args?             │  │
│  │    YES: executeVariadicFn()   │  │  │  │    (controlled by             │  │
│  │    NO:  executeNormalFn()     │  │  │  │     EvaluateArgs flag)        │  │
│  └────────┬──────────────────────┘  │  │  └────────┬──────────────────────┘  │
│           │                         │  │           │                         │
│           v                         │  │           v                         │
│  ┌───────────────────────────────┐  │  │  ┌───────────────────────────────┐  │
│  │ 2. Evaluate Arguments         │  │  │  │ 2. Validate Arg Count         │  │
│  │    - Each arg passed through  │  │  │  │    - Fixed: exact match       │  │
│  │      Evaluate()               │  │  │  │    - Variadic: at least N     │  │
│  └────────┬──────────────────────┘  │  │  └────────┬──────────────────────┘  │
│           │                         │  │           │                         │
│           v                         │  │           v                         │
│  ┌───────────────────────────────┐  │  │  ┌───────────────────────────────┐  │
│  │ 3. Create Child MEM           │  │  │  │ 3. Validate Arg Types         │  │
│  │    - Fork from closure or     │  │  │  │    - Match against            │  │
│  │      current MEM              │  │  │  │      EnvParameter types       │  │
│  │    - Variadic: inject $args   │  │  │  │    - Check :I, :S, :R, etc.   │  │
│  │    - Normal: bind params      │  │  │  └────────┬──────────────────────┘  │
│  └────────┬──────────────────────┘  │  │           │                         │
│           │                         │  │           v                         │
│           v                         │  │  ┌───────────────────────────────┐  │
│  ┌───────────────────────────────┐  │  │  │ 4. Execute Body Function      │  │
│  │ 4. Validate Param Types       │  │  │  │    - Body(ctx, args)          │  │
│  │    - Match against            │  │  │  │    - Direct Go code           │  │
│  │      Parameter.Type           │  │  │  │    - Access Runtime via ctx   │  │
│  └────────┬──────────────────────┘  │  │  └────────┬──────────────────────┘  │
│           │                         │  │           │                         │
│           v                         │  │           v                         │
│  ┌───────────────────────────────┐  │  │  ┌───────────────────────────────┐  │
│  │ 5. Create Child Context       │  │  │  │ 5. Validate Return Type       │  │
│  │    - Same IO, FS              │  │  │  │    - Match against            │  │
│  │    - Same FunctionGroups      │  │  │  │      ReturnType               │  │
│  │    - Child MEM                │  │  │  └────────┬──────────────────────┘  │
│  └────────┬──────────────────────┘  │  │           │                         │
│           │                         │  │           v                         │
│           v                         │  │        RETURN                       │
│  ┌───────────────────────────────┐  │  │                                     │
│  │ 6. Execute Body Instructions  │  │  └─────────────────────────────────────┘
│  │    - Iterate Function.Body    │  │
│  │    - Evaluate each in order   │  │
│  │    - Return last result       │  │
│  └────────┬──────────────────────┘  │
│           │                         │
│           v                         │
│  ┌───────────────────────────────┐  │
│  │ 7. Validate Return Type       │  │
│  │    - Match against            │  │
│  │      Function.ReturnType      │  │
│  └────────┬──────────────────────┘  │
│           │                         │
│           v                         │
│        RETURN                       │
│                                     │
└─────────────────────────────────────┘

EXAMPLE EXECUTION FLOWS:

1. User Function Call:  (my-add 5 10)
   └─> Evaluate(my-add) -> lookup in MEM -> returns OBJ_TYPE_FUNCTION
       └─> executeObjectFunction([5, 10])
           └─> Create child MEM, bind params, execute body, return result

2. Core Function Call:  (putln "hello")
   └─> Evaluate(putln) -> lookup identifier -> returns OBJ_TYPE_IDENTIFIER
       └─> lookupEnvFunction("putln") -> found in "core" FunctionGroup
           └─> executeEnvFunction([evaluated "hello"])
               └─> cmdPutln writes to IO

3. CGS Function Call:  (fs/read_file "test.txt")
   └─> Evaluate(fs/read_file) -> lookup identifier -> returns OBJ_TYPE_IDENTIFIER
       └─> lookupEnvFunction("fs/read_file") -> found in "fs" FunctionGroup
           └─> executeEnvFunction([evaluated "test.txt"])
               └─> cmdReadFile accesses FS interface via Runtime

LEGEND:
  →     Direct transformation
  ↓     Continues to next step
```

## Key Architectural Points

**Function Categories**: The runtime distinguishes between Object Functions (user-defined via `fn`) stored in MEM and Env Functions (runtime-provided) organized into Function Groups. This separation enables controlled extensibility.

**Function Groups**: Each FunctionGroup implements a simple interface exposing a Name() and Functions() map. Core functions live in `env/core.go` while Command Grouped Symbols (CGS) are organized by domain in `pkg/slp/cgs/*`.

**Evaluation Pipeline**: All arguments flow through a validation pipeline that checks count, type, and evaluates based on the function's EvaluateArgs flag. This enables both strict type enforcement and lazy evaluation patterns.

**Memory Scoping**: Object functions capture their defining scope as a closure, forking memory contexts for each invocation. Env functions operate directly within the current evaluation context but can access runtime interfaces (MEM, FS, IO).

**Type System**: Type validation occurs at runtime using type symbols (:I, :S, :R, etc.) with support for :* (any type) wildcards. 

# Tests

The system is reasonably well tested, and all tests can be ran with a simple `make clean && make test`

This will launch all go tests, and it will also cd into the `tests/` directory to launch the series of tests 
that cover the core language and command groups.

While there is decent test coverage there has been no investigation into the memory profile of the runtime. This is all still
very much under development.

## Primary SLPX tests

The "core" language coverage is done in the following files: 

```bash
find tests/primitive -name "*.slpx" -type f -exec wc -lw {} + 
     640    2235 tests/primitive/str.slpx
     842    3201 tests/primitive/list.slpx
     556    1704 tests/primitive/fs.slpx
     392    1481 tests/primitive/reflection.slpx
     569    1956 tests/primitive/bootstrap.slpx
     849    2873 tests/primitive/bits.slpx
      73     436 tests/primitive/main.slpx
     346    1345 tests/primitive/match.slpx
    1240    4209 tests/primitive/numbers.slpx
    5507   19440 total
```

the `main.slpx` file is loaded by the `run.sh` in `tests/` which then begins to `use` each slpx file to initiate tests.

## Primitive Test Process

My goal here was to use the most basic level functionality offered by core.go
to test itself, attempting to detect errors in the first ways it might fail if
something were tampered with in the core that could produce an edge case.

I do this by using lambdas to produce integers by-way of a conditional on
raw-typed objects "0" and "1". The logic internally to yield a 1 or 0 from
a conditional, from within a lambda, producing a checked-for "integer" type
exercises go code that is sensitive to change (for obvious reasons - it impacts
every aspect of calling, and conditions.)

We then check values against lambdas for asserting truthy/falsy values, and
kill execution if our expectations are not met.

In the bootstrap file, what is mentioned above is implemented immediatly, followed by
a long series of similar self<->referencing checks using all commands in, and only 
commands from, core.go.

If bootstrap passes, we can then safely assume that the core of the language is working
and can then proceed to load the tests of commands added via CommandGroups which are part
of the main language expressions, but are logically grouped seperate from the base functions required
to build the language itself, and other groups of main expressions.

They are tested in the order they are needed to be used to test others, as we pollute
the environment with all commands in top level statements on-use. 
We can use functions defined in the files to test others, but most importantly,
we gain the trust that commands tested are commands we can rely on to actually do
the test checking, meaining that for instance "reflection" can rely on "numbers" 
to be working.
