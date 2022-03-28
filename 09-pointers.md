# Pointers

## Creation and Dereferencing

- Type of a pointer to `int`: `*int`
- Address of variable `v`: `&v`
- Dereferencing (getting/setting the memory value) from a pointer stored in variable `p`: `*p`

```go
// Printing the memory address from variable `v`
var v int = 42
fmt.Println(&v)

// Declaring a pointer to an `int` and assigning address of `v`
var p *int = &v

// Printing the value that `p` points to
fmt.Println(*p) // "42"

// Changing the value of `v` using `p`
*p = 27
fmt.Println(v) // "27"
```

- Pointer arithmetic is **not allowed**. So `g := &e[1] - 1` will not work.

```go
e := [3]int8{1, 2, 3}
f := &e[0]
g := &e[1]
fmt.Println(e, f, g)
// [1 2 3] 0xc00009c1d0 0xc00009c1d1
```

## Unsafe stuff

- Using the `unsafe` package, it is actually possible to do pointer arithmetic. In the example below, `h` points to `f`.

```go
e := [3]int8{1, 2, 3}
f := &e[0]
g := &e[1]
int8SizeBytes = int(unsafe.Sizeof(&e[1])) / 8
h := (*int8)(unsafe.Add(unsafe.Pointer(&e[1]), -1*int8SizeBytes))
fmt.Println(e, f, g, h)
// [1 2 3] 0xc00009c1d0 0xc00009c1d1 0xc00009c1d0
```

## Structs and pointers

```go
type myStruct struct {
   foo int
}
var ms *myStruct
fmt.Println(ms) // <nil>
ms = new(myStruct)
fmt.Println(ms) // &{0}
(*ms).foo = 42
fmt.Println(ms)        // &{42}
fmt.Println((*ms))     // {42}
fmt.Println((*ms).foo) // 42
fmt.Println(ms.foo)    // 42 !!! (syntax sugar)
```

- Before initialising a pointer variable, it is set to `nil`, which is the "zero value" for pointers.
- Once we use `new()` to create an empty `struct` inside it, it does have the "zero value" for the struct.
- Go displays pointers with a `&` symbol when printing a `struct` pointer as well as the value of the struct itself.
- Accessing the value of a `struct` field from its pointer would require using parenthesis `(*ms).foo`, but Go has some syntax sugar that makes it possible to simply write `ms.foo`.
- **In functions where values are passed by reference, we usually need to explicitly check for `nil` values.**

## Reference and Value types

- `slice` and `map` are reference types, whereas `array` and `struct` are value types. **This means that you must be careful when operating with `slices` and `maps` as to who has access to the information.** `structs` and `arrays` only have this problem when you're explicitly using pointers.
- This means that whenever you're passing a `slice` or `map` inside your application, the data is always going to be the same, it will not be copied across, but rather passed by reference.

```go
// Array: Original doesn't change.
v := [3]int{1, 2, 3}
w := v
fmt.Println(v, w)
w[1] = 42
fmt.Println(v, w)

// Slice: Both change!
x := []int{1, 2, 3}
y := x
fmt.Println(x, y)
y[1] = 42
fmt.Println(x, y)

// Struct: Original doesn't change.
s := myStruct{foo: 42}
t := s
fmt.Println(s, t)
t.foo = 13
fmt.Println(s, t)

// Dict: Both change!
z := map[string]string{"foo": "bar", "fizz": "buzz"}
u := z
fmt.Println(z, u)
z["foo"] = "zee"
fmt.Println(z, u)
```
