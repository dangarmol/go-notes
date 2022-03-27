# Arrays and Slices

## Arrays

### Creating Arrays

```go
var emptyArray [3]int8               // [0 0 0]
numbers := [3]int8{1, 2, 3}          // Array with a length of 3
autoSize := [...]int8{1, 2, 3, 4, 5} // Array with a length of 5
fmt.Println(numbers)
fmt.Println(autoSize)
```

### Multidimensional Arrays

```go
var identityMatrix [3][3]int = [3][3]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
anotherWay := [3][3]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
```

### Using Arrays

- An array's size must be known at compile time, so it can't change during execution.
- Arrays are always mutable.
- Arrays are processed as values, not as pointers. Assigning an array to a variable copies the whole contents. This can be overriden by using the `&` symbol.

```go
a := [...]int8{1, 2, 3}
b := a
b[1] = 5
fmt.Println(a, b)
c := &a
c[1] = 5
fmt.Println(a, c)
```

## Slices

### Creating Slices

```go
// No need to specify size
a := []int8{1, 2, 3}
```

### Using Slices

- They work similarly to arrays most of the time, but there are some major differences.
- Slices are processed as pointers, so assigning a slice to a new variable, simply assigns the pointer.
- There are certain special functions that work on them, such as `len()` and `cap()` for length and capacity respectively.
- The capacity does not necessarily match the length, it is just the size of the underlying array that is automatically resized by Go.

```go
a := []int8{1, 2, 3}
b := a
b[1] = 5
fmt.Println(a, b)
c := &a
(*c)[1] = 5
fmt.Println(a, c)
```

Slicing arrays or slices works both on arrays and slices.

```go
d := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
e := d[:]   // All elements
f := d[3:]  // 4th element to end
g := d[:6]  // First 6 elements
h := d[3:6] // 4th to 6th elements
fmt.Println(e, f, g, h)
d[5] = 42  // These are still pointers, so all references are updated!
fmt.Println(e, f, g, h)
```

`make()` can also be used to declare slices more explicitly. This can be useful when using `append()` on smaller arrays, as this will avoid resizing the underlying array too often.

```go
i := make([]int8, 3, 100) // Length is 3, but capacity of underlying array is 100.
```

`append()` to a slice:

```go
j := []int8{}
j = append(a, 1)
j = append(a, 2)
j = append(a, 3)
j = append(a, 4, 5, 6) // [1, 2, 3, 4, 5, 6]
```

There is something similar to the Javascript spread operator as well.

```go
k := []int8{1, 2, 3}
l := []int8{4, 5, 6}
k = append(k, l...)
```

Pop elements from front or back:

```go
k = k[1:]        // Remove first element
k = k[:len(k)-1] // Remove last element
```

Remove element from slice:

```go
m := []int{1, 2, 3, 4, 5}
m = append(m[:2], m[3:]...) // Remove element with index 2
fmt.Println(m)

// Careful, assigning a slice to a different variable can cause some weird stuff to happen!
n := []int{1, 2, 3, 4, 5}
o := append(n[:2], n[3:]...) // Remove element with index 2
fmt.Println(o) // [1 2 4 5]
fmt.Println(n) // [1 2 4 5 5] !!!!
```
