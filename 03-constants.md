# Constants

- Collection types (arrays, for example) cannot be mutable, they are inherently variable.
- Constants cannot be shadowed.

## Naming Convention

```go
// No caps used for constants, otherwise they would be exported!
const myConst int = 42
```

## Typed/Untyped Constants

```go
const s float64 = math.Sin(1.57) // This will fail, since it needs to be computed in runtime.
const oneConst int = 42
const anotherConst = 33 // Types can be inferred...
const oops := 23        // ... but this will not work
```

## Enumerated Constants

- `iota` is a special keyword for this and has a few interesting use cases.
- Most commonly used at the package level, outside functions.

```go
const (
   ax = iota
   bx
   cx
)
fmt.Println(ax, bx, cx) // "0 1 2"
```

Practical example:

```go
const (
   errorSpecialist = iota // Use first value for error, since it's the default "zero" value, this avoid mistakes.
   catSpecialist
   dogSpecialist
   snakeSpecialist
)
// Another common option is to use "_ = iota" to throw away the value of the first element.

var specialistType int = catSpecialist
fmt.Printf("%v\n", specialistType == catSpecialist)
```

## Enumeration Expressions

`iota` can have more complex but useful uses.

Bitwise operations:

```go
const (
   _  = iota // Ignore first value
   KB = 1 << (10 * iota)
   MB
   GB
   TB
   PB
   EB
   ZB
   YB
)
fileSize := 40000000000.
fmt.Printf("%.2fGB\n", fileSize/GB)
```

Bitshifting operations:

```go
const (
   isAdmin = 1 << iota // 1
   isHeadquarters      // 10
   canSeeFinancials    // 100

   canSeeAfrica        // 1000
   canSeeAsia          // 10000
   canSeeEurope        // 100000
   canSeeNorthAmerica  // 1000000
   canSeeSouthAmerica  // 10000000
)
var roles byte = isAdmin | canSeeFinancials | canSeeEurope
fmt.Printf("%b\n", roles)
fmt.Printf("Is Admin? %v\n", isAdmin&roles == isAdmin)
fmt.Printf("Is HQ? %v\n", isHeadquarters&roles == isHeadquarters)
```
