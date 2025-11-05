# slpx

expressing something ive been thinking about

## Type Symbols

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
