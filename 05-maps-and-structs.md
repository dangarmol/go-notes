# Maps and Struct

## Maps

### Creating and manipulating Maps

```go
populationDensity := map[string]int{
   "London": 200,
   "Madrid": 100,
   "Zürich": 50,
   "Munich": 40,
}
fmt.Println(populationDensity)
```

- Arrays can be keys for a map, but slices cannot.
- The order of a map is not guaranteed.
- Maps are copied by reference, so assigning it to a different variable and changing it also alters the original.

Empty maps can be created with `make()`

```go
make(map[string]int)
```

Manipulating maps:

```go
populationDensity["Madrid"] = 101
fmt.Println(populationDensity["Madrid"])
populationDensity["Berlin"] = 150 // Add a new key-value pair
delete(populationDensity, "Munich") // Delete a key from the map
fmt.Println(populationDensity["Munich"]) // Prints 0!!
population, ok := populationDensity["Munich"]
fmt.Println(population, ok) // Check if key exists first!
_, ok = populationDensity["Munich"] // More common usage
len(populationDensity) // Returns the length
```

## Structs

- Structs allow us to store any kind of structured data inside arrays, slices or maps.

### Creating Structs

```go
type Doctor struct {
   number   int
   name     string
   patients []string
}

// Specifying keys when using a struct is optional, but recommended.
aDoctor := Doctor{
   Number: 1,
   Name:   "John Smith",
   Patients: []string{
      "Liz",
      "Jo",
      "Mary",
   },
}

fmt.Println(aDoctor.Patients[1])
```

- Passing structs around makes copies of them, not just pointers.
- However, copying by reference is also possible.

```go
aPerson := struct{ name string }{name: "Daniel"}
anotherPerson := aPerson
anotherPerson.name = "Jeff"
fmt.Println(aPerson) // Daniel
fmt.Println(anotherPerson) // Jeff
oneMorePerson := &anotherPerson // References can be made
oneMorePerson.name = "Mike"
fmt.Println(anotherPerson) // Mike
fmt.Println(oneMorePerson) // Mike
```

### Anonymous Structs

These are only useful in limited situations.

```go
aPerson := struct {
   name string
   age  int
}{
   name: "Daniel",
   age:  25,
}
```

### Naming Conventions

- As with every other variable, if they start with a Capital letter, they will be exported, otherwise they won't.

### Embedding (Inheritance)

- Go has no traditional inheritance model, instead, we can embed structs into structs.
- Embedding is not the same as inheritance. According to the example below, a `Bird` is not an `Animal`, but it does contain the same properties, it's just syntactic sugar. To make it interchangeable, we would have to use **interfaces** instead.
- Generally it's much better to use interfaces for common behaviour.

```go
type Animal struct {
   Name   string
   Origin string
}

type Bird struct {
   Animal   // Embedding, means Bird is an animal
   SpeedKPH float32
   CanFly   bool
}

birdie := Bird{}
birdie.Name = "Emu"
birdie.Origin = "Australia"
birdie.SpeedKPH = 48
birdie.CanFly = false
fmt.Println(birdie)

// Using explicit syntax is slightly more complicated
falconBird := Bird{
   Animal:   Animal{Name: "Falcon", Origin: "Spain"},
   SpeedKPH: 200,
   CanFly:   true,
}
fmt.Println(falconBird)
```

### Tags

- They can be useful for validation frameworks, for example.
- They can be obtained using the `reflect` package.

```go
import (
   "fmt"
   "reflect"
)

type Car struct {
   Make   string `required max:"100"`
   Model  string
   Colour string
}

t := reflect.TypeOf(Car{})
field, _ := t.FieldByName("Make")
fmt.Println(field.Tag)
```
