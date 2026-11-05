# dxpr

A dice expression evaluator for the command line.

## Installation

```shell
go install github.com/shvixxl/dxpr/cmd/dxpr@latest
```

## Usage

```shell
dxpr <expression>
```

```shell
dxpr 'd20 + 5'
```

## Expressions

### Dice

```shell
dxpr d20
dxpr 3d6
dxpr d100
dxpr 3d20kh2
dxpr 4d6dl
dxpr 'max(d6, 2)d6'
dxpr '(d10, 5)kh1'
```

### Operators

```shell
dxpr 'd20 + 10'
dxpr '4d8 - 2'
dxpr '2d6 * 3'
dxpr '2d6 x 3'
dxpr 'd100 / 7'
dxpr '2^10'
dxpr '3 ** 2'
dxpr '(d20 + 5) * 2'
```

### Functions

```shell
dxpr 'abs(-5)'
dxpr 'round(7 / 2)'
dxpr 'floor(d20 / 3)'
dxpr 'min(3d6, 10)'
dxpr 'max(d20, 10)'
```
