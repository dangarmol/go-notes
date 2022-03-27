# Primitives

For all the types below, uninitialised variables have a "zero" value, or a default.

- bool = false
- int = 0
- float = 0.0
- string = ""

## Boolean

```go
var yes bool = true
var no bool = false
maybe := 1 == 2 // false
```

### Logical Operators

- `&&` AND
- `||` OR
- `!` NOT

## Numeric

### Numeric Operators

Standard: `+ - * / %`: Types do not change when using them (i.e. `10 / 3 = 3`).

### Bit Operators

- `&` AND
- `|` OR
- `^` XOR
- `&^` AND-NOT
- `<<` LEFT-BIT-SHIFT
- `>>` RIGHT-BIT-SHIFT

### Integers

- `int` is an integer of an unspecified size, depends on the platform, but least 32 bits.
- `int8`, `int16`, `int32`, `int64`: Fixed size signed integers.
- `uint8`, `uint16`, `uint32`: Fixed size unsigned integers.
- There is a `byte` type, which is just an alias for `uint8`.
- There is no `uint64`.
- For larger integers, the `math` package can be used.

### Floating point

Float types follow the IEEE-754 standard.

- `float32`, `float64`: No other sizes available.
- Can be declared using scientific notation: `3.14E15` or `3.13e15`.

### Complex numbers

- Unlike most languages, Go offers `complex64` and `complex128` as primitives.
- There are built-in functions to get the real and imaginary parts.

```go
var n complex64 = 1 + 2i
n = complex(3, 5)
fmt.Println(n, real(n), imag(n))
```

## Other types

### Strings

- Only work for UTF-8 text. They are just an array of `byte`.
- Strings are immutable.
- Nevertheless, string concatenation is allowed using `+`.
- Strings can be subscripted, but the content returned will be a UTF-8 integer code.
- Strings can be declared using **"double quotes"**.

```go
s := "This is a string."
fmt.Println(s[2], string(s[2])) // "105 i"
```

### Bytes

- Many functions work with byte slices (`[]byte`), which makes them more generic. They are just an array of bytes, so they work for strings (array of `uint8` a.k.a. `byte`), as well as files, for instance.
- A byte can be declared using **'single quotes'**.

```go
// Converting a string to []byte
s := "This is a string."
b := []byte(s) // [84 104 105 115 32 105 115 32 97 32 115 116 114 105 110 103 46]
```

### Runes

- Represents UTF-32 characters, using 32-bit integer values.
- `rune` is just an alias for `int32`.
- UTF-8 is valid UTF-32, but UTF-32 does not necessarily use 32 bits per character.
- Usually, special methods are required to process them (e.g. `strings.Reader#ReadRune`).
- A rune can be declared using **'single quotes'**.

```go
var r rune = 'a'
fmt.Printf("%v, %T\n", r, r) // 97, int32
```
