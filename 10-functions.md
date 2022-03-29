# Functions

## Syntax

- The entry point of any application must be a `func main()` function, and it must be within the `package main`.

## Parameters

```go
sayGreeting("Hello", "Daniel")

func sayGreeting(greeting string, name string) {
   fmt.Println(greeting, name)
}

// If all arguments have the same type, it can be put only once
func sayGreeting(greeting, name string) {
   fmt.Println(greeting, name)
}
```

- Passing pointers instead of values allows the function to modify the contents of the variables.
- Additionally, passing pointers to large data structures is a lot more efficient because it doesn't require the data to be copied across.
- When using `slices` or `maps`, the pointers are passed automatically!

```go
func sayGreetingPointer(greeting string, name *string) {
   fmt.Println(greeting, *name)
   *name = "Dani"
}

name := "Daniel"
sayGreetingPointer("Hello", &name)
fmt.Println(name) // "Dani"
```

### Variatic parameters

We can request any number of parameters as the last parameter of a function. The function still allows for other arguments, and the variatic parameters have to be the same type. **There can only be one of these.**

```go
func sum(msg string, values ...int) int {
   fmt.Println(msg, values)
   result := 0
   for _, v := range values {
      result += v
   }
   return result
}
```

## Return values

Some rare features of Go functions:

- Can return a local variable by reference. When Go recognises you want to do this, it stores the variable in the `heap` instead of the `stack`, so that the memory is not deallocated after the function returns.

```go
func sumRef(values ...int) *int {
   result := 0
   for _, v := range values {
      result += v
   }
   return &result
}

fmt.Println(sumRef(1, 2, 3))  // 0xc000282210
fmt.Println(*sumRef(1, 2, 3)) // 6
```

- Named return values: A variable is declared for you and then returned automatically. It is available within the scope of the function.

```go
func sumNamedReturn(values ...int) (result int) {
   for _, v := range values {
      result += v
   }
   return
}
```

- Multiple return values: Quite useful in Go, since Go reserves `panic` calls to severe cases where the app cannot continue, whereas in most cases, there is a better course of action. For example, in a division by zero.

```go
func divide(a, b float64) (float64, error) {
   if b == 0 {
      return 0.0, fmt.Errorf("Cannot divide by zero!")
   }
   return a/b, nil
}

d, err = divide(10.0, 0.0)
if err != nil {
   fmt.Println(err)
   return
}
fmt.Println(d)
```

## Anonymous functions

Although anonymous functions have access to all the variables of the scope where they are declared, usually, it is a better idea to pass the values they need as arguments. This way, we avoid problems when working asynchronously.

```go
func() {
   fmt.Println("Hello Go!")
}()

for i := 0; i < 5; i++ {
   func(i int) {
      fmt.Println(i)
   }(i)
}
```

## Functions as types

Similar to Python's `lambda` functions, we can store functions in variables in Go.

```go
f := func() {
   fmt.Println("Hello Go!")
}
f()
```

The type of functions is `func()`. They can also be declared explicitly as follows:

```go
var f func() = func() {
   // Do something.
}
```

## Methods

- These are functions that are called on an object to perform an action on itself. They are like class methods in other languages.
- To declare such a function, we just need to use the following syntax: `func (x type) funcName {}` or `func (x *type) funcName {}`.
- They can't be declared inside a function, since they're a non-anonymous function.
- By default, the function does create a copy of the object (`value` receiver), so it can't modify it.
- However, we can also make it a `pointer` receiver to be able to change the original object.

```go
type greeter struct {
   greeting string
   name     string
}

func (g greeter) greet() {
   fmt.Println(g.greeting, g.name)
   g.name = "Dan" // This won't have any effect.
}

func (g *greeter) greetAndChange() {
   fmt.Println(g.greeting, g.name)
   g.name = "Dan"
}

func main() {
   g := greeter{
      greeting: "Hello",
      name:     "Daniel",
   }
   g.greet()
   fmt.Println(g.name)
   g.greetAndChange()
   fmt.Println(g.name)
}
```
