# Control Flow

## If

- With the initialiser syntax, the first part of the statement populates variables and the second part after the semicolon `;` evaluates. The variables from the first part of the statement can be used inside the block.

```go
// Standard syntax
if true {
   fmt.Println("Hello there!")
}

// Initialiser syntax
if population, ok := statePopulations["Florida"]; ok {
   fmt.Println(population)
}
```

### Comparison operators

- The following operators are valid in Go: `< <= => > == !=`.
- They can be combined or inverted by using: `&& || !`.
- Additionally, Go uses shortcircuiting:

```go
guess := -5
if guess < 1 || returnTrue() || guess > 100 {
   // returnTrue() will not be executed, since the first element evaluated to `true` already
}
```

- Comparing `float` values can be dangerous, since they are approximations. To use tolerances, the following can be used:

```go
if math.Abs(myNum / myOtherNum - 1) < 0.001 {}
```

### If-else and if-else if

```go
if guess > 0 && guess <= 50 {
   // Do something
} else if guess < 50 && guess < 100 {
   // Do something
} else {
   // Do something
}
```

## Switch

### Simple cases

```go
n := 2
switch n {
   case 1:
      fmt.Println("One")
   case 2:
      fmt.Println("Two")
   default:
      fmt.Println("Not one or two")
}
```

### Complex cases

- Using this syntax, cases cannot be overlapping.
- Initialisers also allowed.

```go
switch n := 3+2; n { // Initialiser also allowed
   case 1, 2, 3:
      fmt.Println("Between one and three")
   case 4, 5, 6:
      fmt.Println("Between four and six")
   default:
      fmt.Println("Not between one and six")
}
```

- Switches can also omit the tag altogether and have "unrelated" comparisons, though the switch will exit as soon as one of them evaluates to `true`.

```go
i := 10
switch { // Empty tag also allowed
   case i >= 1 && i <= 3:
      fmt.Println("Between one and three")
   case i >= 4 && i <= 6:
      fmt.Println("Between four and six")
   default:
      fmt.Println("Not between one and six")
}
```

### Falling through

- To avoid common errors, Go stops after the first statement that evaluates to `true` by default unlike other languages that require `break`, but this can be overriden with the keyword `fallthrough`.

```go
i := 10
switch {
   case i <= 10:
      fmt.Println("Less than or equal to ten")
      fallthrough
   case i <= 20:
      fmt.Println("Less than or equal to twenty")
   default:
      fmt.Println("Greater than twenty")
}
```

### Type switches

```go
var j interface{} = 1
switch j.(type) {
case int:
   fmt.Println("j is an int")
case float64:
   fmt.Println("j is a float64")
case string:
   fmt.Println("j is a string")
case [2]int: // Arrays with different lengths are different types
   fmt.Println("j is a 2-integer array")
default:
   fmt.Println("j is another type")
}
```
