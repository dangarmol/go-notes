# Interfaces

## Basics

- `structs` describe data, whereas `interfaces` describe behaviour.
- For this reason, inside an interface we can define methods rather than variables.
- As a rule of thumb, if an interface has a single function, the name of the interface should be the function + er (Write -> Writer)

```go
type Writer interface {
   Write([]byte) (int, error)
}

type ConsoleWriter struct{}

type FunConsoleWriter struct{}

func (cw ConsoleWriter) Write(data []byte) (int, error) {
   n, err := fmt.Println(string(data))
   return n, err
}

func (cw FunConsoleWriter) Write(data []byte) (int, error) {
   n, err := fmt.Println("Fun", string(data))
   return n, err
}

func interfacesDemo() {
   var cw Writer = ConsoleWriter{}
   cw.Write([]byte("Hello Go!"))
   var fcw Writer = FunConsoleWriter{}
   fcw.Write([]byte("Hello Go!"))
}
```

- In the case above, we have specified the behaviours that a `Writer` must have (`Write`). But many specialised versions of a `Writer` can be implemented.
- Creating the struct and the function below is like the keyword `implements` in other languages. It specifies a particular case of a `Writer` (`ConsoleWriter` in our case) and what its behaviour should be.
- There can be multiple implementations of a `Writer`. A good example of this is a database. There can be multiple types of databases, but usually all of them can perform the same actions (read, write, delete, etc.). In this case, an interface `database` could be created with these actions and there could be different implementations of the same operations depending on the database it's acting upon. Then it becomes just as easy as replacing the initialisation by whichever kind of database is required.
- Any type that can have functions associated can be used.

```go
type Incrementer interface {
   Increment() int
}

type IntCounter int

func (ic *IntCounter) Increment() int {
   *ic++
   return int(*ic)
}

func moreInterfacesDemo() {
   var ic IntCounter = IntCounter(0)
   for i := 0; i < 5; i++ {
      fmt.Println(ic.Increment())
   }
}
```

## Composing interfaces

- It is done exactly the same way as when embedding a `struct` into another one, but just doing the same with `interfaces`.
- This way we can bundle different `interfaces` together, which allows for composition of complex interfaces from several more atomic interfaces.
- **See example on separate file.**

## Type conversion

- We can convert interfaces from one type to another as long as they implement the same methods.
- In this case, when the `wc` variable is a `WriterCloser`, we can access the functions it implements, but we cannot access any of it attributes. By converting it to a `BufferedWriterCloser`, we are allowed to see the internal attributes, in this case, we could access the `buffer` directly by using `bwc.buffer`.
- **See example on separate file.**

```go
// `wc` is actually of type `*BufferedWriterCloser`.
var wc WriterCloser = NewBufferedWriterCloser()
wc.Write([]byte("Hello there, this is another test."))
wc.Close()

// Type conversion from `WriterCloser` to `*BufferedWriterCloser`.
bwc := wc.(*BufferedWriterCloser)
bwc.Write([]byte("We are now trying type conversion."))
bwc.Write([]byte("This is now a *BufferedWriterCloser."))
bwc.Close()

// Conversion from `WriterCloser` to `io.Reader`
ior := wc.(io.Reader) // This will cause a `panic`!
//panic: interface conversion: *main.BufferedWriterCloser is not io.Reader: missing method Read
```

- Recovering from a `panic` is an expensive operation in Go, so there is a better way to check whether a type conversion works with the following syntax:

```go
r, ok := wc.(io.Reader)
if ok {
   fmt.Println(r)
} else {
   fmt.Println("Conversion failed")
}
```

### The empty interface

- An empty interface (defined by type `interface{}`) has no functions associated to it, so it can be used to make a composition of any type, even primitives.
- This is useful for dealing with things that can be multiple types.
- Creating a variable of type `interface{}`, renders it useless as an interface, because it has no methods that can be used. We would need to either do type conversions, or to use the `reflect` package in order to figure out what kind of object we are dealing with.
- Empty interfaces are quite common, but they are almost always an intermediate step.

```go
var myObj interface{} = NewBufferedWriterCloser()
if wc, ok := myObj.(WriterCloser); ok {
   wc.Write([]byte("Just one more test using type conversion "))
   wc.Write([]byte("onto an empty interface."))
   wc.Close()
} else {
   fmt.Println("Conversion failed")
}
```

### Type switches

- Because empty interfaces have no need to implement any functions, they can be used to store any type of data.

```go
var i interface{} = 0
switch i.(type) {
case int:
   fmt.Println("'i' is an integer")
case string:
   fmt.Println("'i' is a string")
default:
   fmt.Println("'i' is something else")
}
```

## Implementing with values vs pointers

- If any of the methods require a pointer receiver, the whole interface will have to be passed as a pointer, but not the other way around.
- **Pointer receivers are required when the internal value of the struct wants to be persisted across calls!**

```go
//////////////////////////////////
// Implementation using a value //
//////////////////////////////////

type ValueWriterCloser struct {
   saved string
}

func (vwc ValueWriterCloser) Write(data []byte) (int, error) {
   fmt.Println("Current saved data:", vwc.saved)
   vwc.saved = string(data)
   return len(vwc.saved), nil
}

func (vwc ValueWriterCloser) Close() error {
   return nil
}

////////////////////////////////////
// Implementation using a pointer //
////////////////////////////////////

type PointerWriterCloser struct {
   saved string
}

func (pwc *PointerWriterCloser) Write(data []byte) (int, error) {
   fmt.Println("Current saved data:", pwc.saved)
   pwc.saved = string(data)
   return len(pwc.saved), nil
}

func (pwc *PointerWriterCloser) Close() error {
   return nil
}

func valueVsPointerDemo() {
   var vwc WriterCloser = ValueWriterCloser{}
   fmt.Println(vwc)
   var vwc2 WriterCloser = &ValueWriterCloser{}
   fmt.Println(vwc2)
   var pwc WriterCloser = &PointerWriterCloser{}
   fmt.Println(pwc)
}
```

A more visual example:

```go
type Point struct {
   x, y float64
}

// If we make `p` of type `Point` instead of `*Point`,
// changes will not be persisted!
func (p *Point) scale(s float64) {
   p.x = p.x * s
   p.y = p.y * s
}

// `size()` does not need to have a pointer receiver,
// since it does not modify the data inside `p`. However, it will
// make a copy of `p` in memory.
func (p Point) size() float64 {
   return math.Sqrt(p.x * p.x + p.y * p.y)
}

func main() {
   p := &Point{3, 4}
   fmt.Printf("Before scaling: %+v, size: %v\n)", p, p.size())
   // "Before scaling: &{x:3 y:4}, size: 5"
   var sc float64 = 5
   p.scale(sc)
   fmt.Printf("After scaling: %+v, size: %v\n)", p, p.size())
   // "After scaling: &{x:15 y:20}, size: 25"
}
```

## Best practices

- Use many, small interfaces.
  - Single method interfaces are some of the most powerful and flexible: `io.Writer`, `io.Reader`, `interface{}`.
- Don't export interfaces for types that will be consumed.
  - Often, it is perfectly acceptable to export the concrete type. As an example in the [database/sql](https://pkg.go.dev/database/sql) package, the `DB` type is exported as a concrete struct, and it has all sorts of methods that pointed to other concrete structs.
  - By not exporting the interface, it allows the consumer to define their own interface that can be used for testing, for example, without having to implement every method.
  - In the case of `DB`, if you don't need every method in the `DB` type, your own interface doesn't need to have all of them, only the ones you require.
- Do export interfaces for types that will be used by package.
  - If you will be accepting a value, try and accept an interface, if at all possible.
- Design functions and methods to receive interfaces whenever possible.
  - If you're just accepting data, this is often not possible, but if you're taking behaviour providers as arguments, it is often useful to accept interfaces instead.
