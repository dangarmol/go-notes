# Defer, Panic and Recover

## Defer

Delays the execution of any functions passed to it until the final statement is executed. It executes the statements just before returning.

```go
// Prints: start, end, middle
fmt.Println("start")
defer fmt.Println("middle")
fmt.Println("end")

// Prints: end, middle, start (LIFO)
defer fmt.Println("start")
defer fmt.Println("middle")
defer fmt.Println("end")
```

A useful pattern in which `defer` is used is by opening a resource, checking for errors and immediately doing a deferred close. This way, the close will always be done and it comes together as a block. Since it works in a LIFO fashion, it will close resources that generated afterwards before the older ones.

```go
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
```

This can be problematic when deferring closing resources inside a loop, since there can be many resources open at once and they will only close once the function ends.

```go
a := "start"
defer fmt.Println(a) // Prints "start". Takes the arguments when called, not when executed.
a = "end
```

## Panic (Exceptions)

In most languages, there are things that would be considered exceptions, but in go they aren't. For example, opening a file that doesn't exist. Nevertheless, if for any reason there is a situation in which a program cannot continue, that is what Go calls `panic`.

```go
// This will stop the execution immediately.
panic("Something nasty happened here!")
```

Most functions in Go will not declare `panic` themselves, but rather they will return an `error` value and it's up to the developer whether to declare a `panic` or not.

The following snippet is a good example of when to declare a panic. The execution will be interrupted explicitly by the developer in this basic web server if there is already an application using port 8080 because of the call to `panic`.

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
   w.Write([]byte("Hello Go!"))
})
err := http.ListenAndServe(":8080", nil)
if err != nil {
   panic(err.Error())
}
```

## Panic + Defer

Another reason why `defer` is useful, is because it can **clean up** in the event of a `panic` during the execution. The deferred function will attempt to run before stopping the execution and showing the user the `panic`.

```go
fmt.Println("start")
defer fmt.Println("deferred")
if doIt {
   panic("Something nasty happened here!")
}
fmt.Println("end")
```

## Recover (try-catch)

When you `recover()` a `panic`, the function in which the `panic` happened, does not continue execution, but the rest of the program continues as usual. An example of this is the following:

```go
func main() {
   fmt.Println("start")
   panicker()
   fmt.Println("end")
}
func panicker() {
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
```

- `recover()` simply returns whether there is an active `panic` for the developer to try and recover from.
- If nothing handles the panic by calling the `recover()` function, then the program will exit.
- It is only useful in deferred functions, otherwise it will not be executed.
