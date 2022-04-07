package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"unsafe"
)

var (
	hello int = 1
	hi    int = 2
)

var j float32 = 0.5

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func primitivesDemo() {
	var i int = 25
	i = 42
	fmt.Println("Hello Go!")
	fmt.Println(hello, hi)
	fmt.Printf("%v, %T\n", i, i)
	fmt.Printf("%v, %T\n", j, j)
	fmt.Println(i, strconv.Itoa(i)) // string(i) returns "*"

	var yes bool = true
	var no bool = false
	maybe := 1 == 2
	fmt.Println(yes, no, maybe)

	var n complex64 = 1 + 2i
	n = complex(3, 5)
	fmt.Println(n, real(n), imag(n))

	s := "This is a string."
	fmt.Println(s[2], string(s[2])) // "105 i"
	b := []byte(s)
	fmt.Println(b) // [84 104 105 115 32 105 115 32 97 32 115 116 114 105 110 103 46]

	var r rune = 'a'
	fmt.Printf("%v, %T\n", r, r) // 97, int32
}

func enumeratedConstantsDemo() {
	const (
		ax = iota
		bx
		cx
	)
	fmt.Println(ax, bx, cx)

	fileSize := 40000000000.
	fmt.Printf("%.2fGB\n", fileSize/GB)

	const (
		isAdmin          = 1 << iota // 1
		isHeadquarters               // 10
		canSeeFinancials             // 100

		canSeeAfrica       // 1000
		canSeeAsia         // 10000
		canSeeEurope       // 100000
		canSeeNorthAmerica // 1000000
		canSeeSouthAmerica // 10000000
	)
	var roles byte = isAdmin | canSeeFinancials | canSeeEurope
	fmt.Printf("%b\n", roles)
	fmt.Printf("Is Admin? %v\n", isAdmin&roles == isAdmin)
	fmt.Printf("Is HQ? %v\n", isHeadquarters&roles == isHeadquarters)
}

func arraysDemo() {
	fmt.Println("Arrays")
	var emptyArray [3]int8               // [0 0 0]
	numbers := [3]int8{1, 2, 3}          // Array with a length of 3
	autoSize := [...]int8{1, 2, 3, 4, 5} // Array with a length of 5
	fmt.Println(emptyArray)
	fmt.Println(numbers)
	fmt.Println(autoSize)
	var identityMatrix [3][3]int = [3][3]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
	anotherWay := [3][3]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
	fmt.Println(identityMatrix)
	fmt.Println(anotherWay)

	a := [...]int8{1, 2, 3}
	b := a
	b[1] = 5
	fmt.Println(a, b)
	c := &a
	c[1] = 5
	fmt.Println(a, c)

	d := [...]int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	e := d[:]   // all elements
	f := d[3:]  // 4th element to end
	g := d[:6]  // first 6 elements
	h := d[3:6] // 4th to 6th elements
	fmt.Println(e, f, g, h)
	fmt.Printf("%v, %T\n", e, e)
	d[5] = 42 // These are just pointers, so all references are updated!
	fmt.Println(e, f, g, h)
}

func slicesDemo() {
	fmt.Println("Slices")
	a := []int8{1, 2, 3}
	b := a
	b[1] = 5
	fmt.Println(a, b)
	c := &a
	(*c)[1] = 5
	fmt.Println(a, c)

	d := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	e := d[:]   // all elements
	f := d[3:]  // 4th element to end
	g := d[:6]  // first 6 elements
	h := d[3:6] // 4th to 6th elements
	fmt.Println(e, f, g, h)
	fmt.Printf("%v, %T\n", e, e)
	d[5] = 42 // These are still pointers, so all references are updated!
	fmt.Println(e, f, g, h)

	j := []int8{}
	j = append(j, 1)
	j = append(j, 2)
	j = append(j, 3)
	j = append(j, 4, 5, 6)
	fmt.Println(j)

	k := []int8{1, 2, 3}
	l := []int8{4, 5, 6}
	k = append(k, l...)
	fmt.Println(k)
	k = k[1:]        // Remove first element
	k = k[:len(k)-1] // Remove last element
	fmt.Println(k)

	m := []int{1, 2, 3, 4, 5}
	m = append(m[:2], m[3:]...) // Remove element with index 2
	fmt.Println(m)

	n := []int{1, 2, 3, 4, 5}
	o := append(n[:2], n[3:]...) // Remove element with index 2
	fmt.Println(o)               // [1 2 4 5]
	fmt.Println(n)               // [1 2 4 5 5] !!!!
}

func mapsDemo() {
	populationDensity := map[string]int{
		"London": 200,
		"Madrid": 100,
		"Zürich": 50,
		"Munich": 40,
	}
	fmt.Println(populationDensity)
	populationDensity["Madrid"] = 101
	fmt.Println(populationDensity["Madrid"])
	populationDensity["Berlin"] = 150
	delete(populationDensity, "Munich")
	fmt.Println(populationDensity)
	fmt.Println(populationDensity["Munich"]) // Prints 0 !!!
	population, ok := populationDensity["Munich"]
	fmt.Println(population, ok)         // Returns if key is ok as well
	_, ok = populationDensity["Munich"] // More common usage
	fmt.Println(ok)
	fmt.Println(len(populationDensity))
}

func structsDemo() {
	type Doctor struct {
		number   int
		name     string
		patients []string
	}

	aDoctor := Doctor{
		number: 1,
		name:   "John Smith",
		patients: []string{
			"Liz",
			"Jo",
			"Mary",
		},
	}

	fmt.Println(aDoctor)
	fmt.Println(aDoctor.patients[1])

	somePerson := struct {
		name string
		age  int
	}{
		name: "Daniel",
		age:  25,
	}
	fmt.Println(somePerson)

	aPerson := struct{ name string }{name: "Daniel"}
	anotherPerson := aPerson
	anotherPerson.name = "Jeff"
	fmt.Println(aPerson)       // Daniel
	fmt.Println(anotherPerson) // Jeff
	oneMorePerson := &anotherPerson
	oneMorePerson.name = "Mike"
	fmt.Println(anotherPerson) // Mike
	fmt.Println(oneMorePerson) // Mike

	type Animal struct {
		Name   string
		Origin string
	}

	type Bird struct {
		Animal   // Embedding, means Bird is an animal
		SpeedKPH float32
		CanFly   bool
	}

	emuBird := Bird{}
	emuBird.Name = "Emu"
	emuBird.Origin = "Australia"
	emuBird.SpeedKPH = 48
	emuBird.CanFly = false
	fmt.Println(emuBird)

	// Using explicit syntax is slightly more complicated
	falconBird := Bird{
		Animal:   Animal{Name: "Falcon", Origin: "Spain"},
		SpeedKPH: 200,
		CanFly:   true,
	}
	fmt.Println(falconBird)

	type Car struct {
		Make   string `max:"100"`
		Model  string
		Colour string
	}

	t := reflect.TypeOf(Car{})
	field, _ := t.FieldByName("Make")
	fmt.Println(field.Tag)
}

func controlFlowDemo() {
	if true {
		fmt.Println("Hello there!")
	}

	guess := -5
	if guess < 1 || guess > 100 {
		fmt.Println("Outside valid range.")
	}

	n := 2
	switch n {
	case 1:
		fmt.Println("One")
	case 2:
		fmt.Println("Two")
	default:
		fmt.Println("Not one or two")
	}

	switch n := 3 + 2; n { // Initialiser also allowed
	case 1, 2, 3:
		fmt.Println("Between one and three")
	case 4, 5, 6:
		fmt.Println("Between four and six")
	default:
		fmt.Println("Not between one and six")
	}

	i := 10
	switch { // Empty tag also allowed
	case i >= 1 && i <= 3:
		fmt.Println("Between one and three")
	case i >= 4 && i <= 6:
		fmt.Println("Between four and six")
	default:
		fmt.Println("Not between one and six")
	}

	switch {
	case i <= 10:
		fmt.Println("Less than or equal to ten")
		fallthrough // This ignores the logic of the next cases and executes directly !!!
	case i <= 20:
		fmt.Println("Less than or equal to twenty")
	default:
		fmt.Println("Greater than twenty")
	}

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
}

func loopingDemo() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	for i, j := 0, 0; i < 5; i, j = i+1, j+1 {
		fmt.Println(i, j)
	}

	// Initialiser and incrementer are optional.
	i := 0
	for i < 5 {
		fmt.Println(i)
		i++
	}

	i = 0
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

Loop:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fmt.Println(i * j)
			if i*j >= 3 {
				break Loop
			}
		}
	}

	s := []int{10, 20, 30}
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
}

func deferHTTPDemo(doPrint bool) {
	res, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close() // To avoid forgetting about closing we can do it here
	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if doPrint {
		fmt.Printf("%s", robots)
	}
}

func recoverPanicDemo() {
	fmt.Println("About to panic...")
	defer func() { // try...
		if err := recover(); err != nil { //...catch
			log.Println("Error:", err)
			// panic(err) // If you actually can't deal with the issue...
		}
	}() // Anonymous function
	panic("Something nasty happened here!")
	// fmt.Println("Actually it's fine") // This would never run :(
}

func deferPanicRecoverDemo() {
	fmt.Println("start")
	defer fmt.Println("middle")
	fmt.Println("end")

	defer fmt.Println("start")
	defer fmt.Println("middle")
	defer fmt.Println("end")

	deferHTTPDemo(true)
	recoverPanicDemo()
	fmt.Println("I am not panicking, we are good again!")
}

func pointersDemo() {
	a := 42
	b := a
	fmt.Println(a, b)
	a = 27
	fmt.Println(a, b)

	var c int = 42
	var d *int = &c
	fmt.Println(c, d)
	fmt.Println(c, *d)
	c = 14
	fmt.Println(c, *d)
	*d = 27
	fmt.Println(c, *d)

	e := [3]int8{1, 2, 3}
	f := &e[0]
	g := &e[1]
	int8SizeBytes := int(unsafe.Sizeof(&e[1])) / 8
	h := (*int8)(unsafe.Add(unsafe.Pointer(&e[1]), -1*int8SizeBytes))
	fmt.Println(e, f, g, h)

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
}

func sayGreeting(greeting, name string) {
	fmt.Println(greeting, name)
}

func sayGreetingPointer(greeting string, name *string) {
	fmt.Println(greeting, *name)
	*name = "Dani"
}

func sum(msg string, values ...int) int {
	fmt.Println(msg, values)
	result := 0
	for _, v := range values {
		result += v
	}
	return result
}

func sumRef(values ...int) *int {
	result := 0
	for _, v := range values {
		result += v
	}
	return &result
}

func sumNamedReturn(values ...int) (result int) {
	for _, v := range values {
		result += v
	}
	return
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0.0, fmt.Errorf("Cannot divide by zero!")
	}
	return a / b, nil
}

func functionsDemo() {
	name := "Daniel"
	sayGreeting("Hello", name)
	fmt.Println(name)

	sayGreetingPointer("Hello", &name)
	fmt.Println(name)

	fmt.Println(sum("Got values:", 1, 2, 3, 4, 5, 6, 7))

	fmt.Println(sumRef(1, 2, 3))  // 0xc000282210
	fmt.Println(*sumRef(1, 2, 3)) // 6

	fmt.Println(sumNamedReturn(1, 2, 3))

	d, err := divide(10.0, 5.0)
	if err != nil {
		fmt.Println(err)
		// Usually `return` here
	} else {
		fmt.Println(d)
	}

	d, err = divide(10.0, 0.0)
	if err != nil {
		fmt.Println(err)
		// Usually `return` here
	} else {
		fmt.Println(d)
	}

	// Anonymous functions
	for i := 0; i < 5; i++ {
		func(i int) {
			fmt.Println(i)
		}(i)
	}

	f := func() {
		fmt.Println("Hello Go!")
	}
	f()
}

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

func methodsDemo() {
	g := greeter{
		greeting: "Hello",
		name:     "Daniel",
	}
	g.greet()
	fmt.Println(g.name)
	g.greetAndChange()
	fmt.Println(g.name)
}

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

func main() {
	primitivesDemo()
	enumeratedConstantsDemo()
	arraysDemo()
	slicesDemo()
	mapsDemo()
	structsDemo()
	controlFlowDemo()
	loopingDemo()
	deferPanicRecoverDemo()
	pointersDemo()
	functionsDemo()
	methodsDemo()
	interfacesDemo()
	moreInterfacesDemo()
}
