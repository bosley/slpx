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


# Something to clarify in a document:

We use `$error` as an injected identifier inside the body of the handler for `try` statements. After this handling, `$error` is not available. 

The user cannot define `$<IDENTIFIER>` as the `$` symbol

----------------------------

