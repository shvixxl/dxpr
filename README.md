# dxpr

A dice expression evaluator for the command line.

## Usage

```shell
dxpr <expression>
```

All arguments are joined into a single expression, so quotes are optional:

```shell
dxpr d20 + 5
```

## Expressions

### Dice

Roll dice with `NdM` notation — `N` dice with `M` sides each. The results are summed.

```shell
dxpr d20        roll one 20-sided die
dxpr 3d6        roll three 6-sided dice and sum them
dxpr d100       roll one 100-sided die
```

The die count may be omitted and defaults to 1: `d20` is the same as `1d20`.

### Arithmetic

Combine rolls with `+`, `-`, `*`, `/` (integer division), `//` (floor division), `%` (modulus) and `^` or `**` (exponentiation):

```shell
dxpr d20 + 10
dxpr 4d8 - 2
dxpr 2d6 * 3
dxpr d100 / 7
dxpr 2^10
dxpr 3 ** 2
```

`^` and `**` bind tighter than `*` and `/`; `2^3^2` is `2^(3^2)` = 512.

Unary negatives are not supported yet; `d20 - 5` is fine, `d20 + -5` is not.

Division by zero and negative exponents are errors.

## Errors

Invalid expressions exit with a non-zero status and print an error message:

```shell
dxpr "d20 +"
dxpr: expected number, got end of expression
```
