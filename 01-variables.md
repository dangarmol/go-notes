# Variables

## Variable Declaration and Assignment

```go
var i int      // Simply declare variable
var i int = 42 // Declare and assign variable
i := 42        // Infere variable type from value
fmt.Printf("%v, %T\n", i, i)  // Print value and type
```

Variables can be declared outside of functions, but then the type needs to be specified explicitly.

They can also be declared in groups.

```go
var (
   name            string = "Daniel"
   favouriteNumber int = 9
   homeURL         string = "https://dangarmol.dev/"
)
```

## Variable Shadowing

This is valid code:

```go
var i int = 35

func main() {
   var i int = 42
   i = 23
}
```

This is not:

```go
var i int = 35
var i int = 42 // This will fail
i := 23        // This will fail too
```

## Scopes

When declaring a variable at the package level, if its name stars with a lowercase letter, it is exposed within the package only, whereas if it is in uppercase, it is exposed publicly outside the package. **This also applies to functions.**

## Number type conversion

```go
i := 42
f := 4.2

a = float32(i)      // Works fine
b = int(f)          // Works, but loses information
c = string(i)       // This will convert to unicode char 42! -> '*'
d = strconv.Itoa(i) // Import "strconv"
```
