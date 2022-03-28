# Looping

## For statements

They are the only way to loop in Go. No `while` loops exist.

### Simple loops

For loops commonly have 3 statements.

- `Initialiser`: Usually a counter, but can initialise anything.
- `Boolean statement`: Exit condition, evaluates to `bool`.
- `Incrementer`: Action to be performed between iterations of the loop.

```go
for i := 0; i < 5; i++ {
   fmt.Println(i)
}

// for-loops can have multiple variables, but only one statement is allowed
for i, j := 0, 0; i < 5; i, j = i+1, j+1 {
   fmt.Println(i, j)
}

// Initialiser and incrementer are optional. For example:
i := 0
for ; i < 5; {
   fmt.Println(i)
   i++
}

// Or similarly to a while loop in other languages...
i := 0
for i < 5 {
   fmt.Println(i)
   i++
}

// Infinite loops are also allowed with `break` and `continue` keywords...
i := 0
for {
   if i > 5 {
      break
   }
   i++
   if i%2 == 0 {
      continue
   }
   fmt.Println(i)
}
```

### Exiting early

- We can use labels, such as `Loop:` to be able to exit the selected loop. Otherwise, we would only `break` out of the inner loop but continue looping outside.

```go
Loop:
   for i := 1; i <= 3; i++ {
      for j := 1; j <= 3; j++ {
         fmt.Println(i * j)
         if i*j >= 3 {
            break Loop
         }
      }
   }
```

### Looping through collections

- This can be achieved using the `range` keyword. This type of `for` loop, makes both the index and the value available.

```go
s := []int{1, 2, 3}
for k, v := range s {
   fmt.Println(k, v)
}

// This works with any container or iterable
sentence := "Hello World!"
for index, letter := range sentence {
   fmt.Println(index, string(letter))
}

populationDensity := map[string]int{
   "London": 200,
   "Madrid": 100,
   "Zürich": 50,
   "Munich": 40,
}
for key, value := range populationDensity {
   fmt.Println(key, value)
}
for key := range populationDensity {
   fmt.Println(key)
}
for _, value := range populationDensity {
   fmt.Println(value)
}
```
