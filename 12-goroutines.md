# Goroutines

## Basics

Most programming languages use OS threads, but these are heavyweight and take time to initialise, which means that spawning thousands of these threads can be impractical or impossible. For contrast, Go uses `goroutines` (`green threads`) instead, which are much more lightweight and are managed in the `user space`, not the `kernel space` by the Go runtime.
Because it is very expensive to create threads in other programming languages, solutions like `thread pooling` have been developed.
**The Go runtime has a scheduler that maps goroutines to OS threads and assigns a certain amount of time on each available thread.**

## Green Threads

Extracted from [StackOverflow](https://softwareengineering.stackexchange.com/a/222694/293004)

**Green threads are threads that are scheduled by a runtime library or virtual machine (VM) instead of natively by the underlying operating system.** Green threads emulate multithreaded environments without relying on any native OS capabilities, and they are managed in user space instead of kernel space, enabling them to work in environments that do not have native thread support.

Go (or more exactly the two existing implementations) is a language producing native code only - it does not use a VM. Furthermore, the scheduler in the current runtime implementations relies on OS level threads (even when `GOMAXPROCS=1`). So I think talking about green threads for the Go model is a bit abusive.

Go people have coined the goroutine term especially to avoid the confusion with other concurrency mechanisms (such as coroutines or threads or lightweight processes).

Of course, Go supports a M:N threading model, but it looks much closer to the Erlang process model, than to the Java green thread model.

Here are a few benefits of the Go model over green threads (as implemented in early JVM):

- Multiple cores or CPUs can be effectively used, in a transparent way for the developer. **With Go, the developer should take care of concurrency. The Go runtime will take care of parallelism.** Java green threads implementations did not scale over multiple cores or CPUs.
- System and C calls are non blocking for the scheduler (all system calls, not only the ones supporting multiplexed I/Os in event loops). Green threads implementations could block the whole process when a blocking system call was done.
- Copying or segmented stacks. In Go, there is no need to provide a maximum stack size for the goroutine. The stack grows incrementally as needed. A goroutine does not require much memory (4KB-8KB, as opposed to ~1MB per OS thread), so a huge number of them can be happily spawned. Goroutine usage can therefore be pervasive.

Now, to address the criticisms:

- With Go, you do not have to write a user space scheduler: it is already provided with the runtime. It is a complex piece of software, but it is the problem of Go developers, not of Go users. Its usage is transparent for Go users. Among the Go developers, [Dmitri Vyukov](http://www.1024cores.net/) is an expert in lockfree/waitfree programming, and he seems to be especially interested in addressing the eventual performance issues of the scheduler. The current scheduler implementation is not perfect, but it will improve.
- Synchronization brings performance problems and complexity: this is partially true with Go as well. But note the Go model tries to promote the usage of channels and a clean decomposition of the program in concurrent goroutines to limit synchronization complexity (i.e. sharing data by communicating, instead of sharing memory to communicate). By the way, the reference Go implementation provides a number of tools to address performance and concurrency issues, like a [profiler](http://blog.golang.org/profiling-go-programs), and a [race detector](http://blog.golang.org/race-detector).
- Regarding page fault and "multiple threads faking", please note Go can schedule a goroutine over multiple system threads. When one thread is blocked for any reason (page fault, blocking system calls), it does not prevent the other threads to continue to schedule and run other goroutines. Now, it is true that a page fault will block the OS thread, with all the goroutines supposed to be scheduled on this thread. However in practice, the Go heap memory is not supposed to be swapped out. This would be the same in Java: garbage collected languages do not accomodate virtual memory very well anyway. If your program must handle page fault in a graceful way, it is probably because it has to manage some off-heap memory. In that case, wrapping the corresponding code with C accessor functions will simply solve the problem (because again C calls or blocking system calls never block the Go runtime scheduler).

Note that the “green thread” article on Wikipedia has been changed to state “threads that are scheduled by a runtime library or virtual machine (VM)”; which means that by that definition your answer would not be correct anymore, as the **Go runtime does the scheduling/management**. I think it's more helpful to define green threads as user-space threads contrasting OS-threads. In which case, goroutines are green threads for sure.

## Creating Goroutines

This will not print anything. The `goroutine` is created, but not given any time to finish before the program ends execution.

```go
func main() {
   go sayHello()
}

func sayHello() {
   fmt.Println("Hello")
}
```

This option will work, but an arbitrary wait is never a good option. This is just for demonstration purposes.

```go
func main() {
   go sayHello()
   time.Sleep(100 * time.Millisecond)
}

func sayHello() {
   fmt.Println("Hello")
}
```

Using an anonymous function will work as well, but it could have unintended consequences. The fact that the anonymous function is a **closure** means that it has access to the variables in the scope of the main function. The danger being that the variables can change before the goroutine has a chance to execute.

```go
func main() {
   var msg = "Hello"
   go func() {
      fmt.Println(msg) // Will print "Goodbye"
   }()
   go func(msg string) {
      fmt.Println(msg) // Will print "Hello"
   }(msg)
   msg = "Goodbye"
   time.Sleep(100 * time.Millisecond)
}
```

This will work because Go makes sure that the anonymous function has access to the variable `msg`, even though they're running on separate threads. Nevertheless, this is not a good practice because it will cause a **race condition**, since the contents of the variable may change before the anonymous function has a chance to execute.

The anonymous function below takes an argument and copies the value of the variable, so it will print as expected.

One more thing to note is that "Hello" and "Goodbye" may be printed in any order.

## Synchronisation

### WaitGroups

`WaitGroups` are used to synchronise multiple goroutines together. There are 3 methods that can be used for this purpose:

- `wg.Add(delta int)`: Adds the specified amount of goroutines to the `WaitGroup`. This can be increased over time or subtracted using negative values of delta. It specifies how many goroutines should be awaited.
- `wg.Done()`: Must be called when a goroutine finishes. Decreases the `WaitGroup` counter by one.
- `wg.Wait()`: Waits until the `WaitGroup` count is zero.

This helps minimise the amount of time that the program must wait as much as possible.

```go
func waitGroupExample() {
   startTime := time.Now()
   fmt.Println("WaitGroup example:")
   var msg = "Hello"
   wg.Add(2) // How many goroutines should be awaited
   go func() {
      fmt.Println(msg) // Will print "Goodbye"
      wg.Done()
   }()
   go func(msg string) {
      fmt.Println(msg) // Will print "Hello"
      wg.Done()
   }(msg)
   msg = "Goodbye"
   wg.Wait() // Waits until the WaitGroup count is zero.
   fmt.Printf("Time waited: %v\n", time.Since(startTime))
}
```

However, it can be tricky to synchronise tasks using a `WaitGroup`. They simply try to accomplish the task as fast as they can without syncing the threads.

The following example will just print out numbers somewhat randomly depending on what the scheduler decides to allocate higher priority to. Numbers may be repeated and even out of order!

```go
var counter = 0

func sayHello() {
   fmt.Printf("Hello #%v\n", counter)
   wg.Done()
}

func increment() {
   counter++
   wg.Done()
}

func forWaitGroupExample() {
   for i := 0; i < 10; i++ {
      wg.Add(2)
      go sayHello()
      go increment()
   }
   wg.Wait()
}
```

### Mutexes

By using a `mutex`, we can ensure that no operations will happen out of order, which will cause the numbers to only go up in ascending order, but in this example, we cannot ensure that the threads will be synchronised and will print all numbers. Most likely there will be repetitions.

```go
var mtx = sync.RWMutex{}
var counterMutex = 0

func sayHelloMutex() {
   mtx.RLock()
   fmt.Printf("Hello #%v\n", counterMutex)
   mtx.RUnlock()
   wg.Done()
}

func incrementMutex() {
   mtx.Lock()
   counterMutex++
   mtx.Unlock()
   wg.Done()
}

func mutexExample() {
   fmt.Println("Simple RWMutex example:")
   runtime.GOMAXPROCS(100)
   for i := 0; i < 10; i++ {
      wg.Add(2)
      go sayHelloMutex()
      go incrementMutex()
   }
   wg.Wait()
}
```

By placing the locks in the same context, we have achieved synchronisation. Once each lock is placed, it will be unlocked asynchronously whenever the goroutine decides, but there will be no repetitions of reads or writes. **However, by doing this, the parallelism is ruined and the code will run sequentially as if there were no goroutines at all.**

```go
var mtx = sync.RWMutex{}
var counterBetterMutex = 0

func sayHelloBetterMutex() {
   fmt.Printf("Hello #%v\n", counterBetterMutex)
   mtx.RUnlock()
   wg.Done()
}

func incrementBetterMutex() {
   counterBetterMutex++
   mtx.Unlock()
   wg.Done()
}

func betterMutexExample() {
   fmt.Println("Better RWMutex example:")
   runtime.GOMAXPROCS(100)
   for i := 0; i < 10; i++ {
      wg.Add(2)
      mtx.RLock()
      go sayHelloBetterMutex()
      mtx.Lock()
      go incrementBetterMutex()
   }
   wg.Wait()
}
```

## Parallelism with runtime.GOMAXPROCS()

- The number of OS threads available can be queried with `runtime.GOMAXPROCS(-1)`.
- By default, `runtime.GOMAXPROCS(-1)` will return the number of cores from the computer.
- This value can be modified using the same function with a positive number, and multiple Go threads can use a single CPU thread at once.
- This means that, in general, it is a good idea to allocate more threads than CPU cores. However, if this number is too high, the overhead of the scheduler and the memory that each thread requires will start to be greater than the benefit it brings.
- This value will eventually go away once the scheduler becomes good enough. Until then, this value should be fine tuned for each application.
- Finally, if `runtime.GOMAXPROCS(1)` is set to 1, the application will run on a single thread without any operations in parallel. In general this is not recommended, but it can be useful in certain situations.

## Best practices

- Don't create goroutines in libraries. Let the consumer control concurrency. A potential exception is if a `channel` is being returned.
- When creating a goroutine, know how it will end. If you don't know when it will end, there is a chance that it will continue running forever and leak memory.
- **Check for race conditions** at compile time by adding the `-race` flag to the compiler. For example: `go run -race main.go`. Example output:

```go
==================
WARNING: DATA RACE
Write at 0x00c000012030 by main goroutine:
  main.waitGroupExample()
      /Users/daniel/go/src/github.com/dangarmol/go-notes/03-goroutines-examples/main.go:106 +0x251
  main.main()
      /Users/daniel/go/src/github.com/dangarmol/go-notes/03-goroutines-examples/main.go:129 +0xa4

Previous read at 0x00c000012030 by goroutine 10:
  main.waitGroupExample.func1()
      /Users/daniel/go/src/github.com/dangarmol/go-notes/03-goroutines-examples/main.go:99 +0x3a

Goroutine 10 (finished) created at:
  main.waitGroupExample()
      /Users/daniel/go/src/github.com/dangarmol/go-notes/03-goroutines-examples/main.go:98 +0x170
  main.main()
      /Users/daniel/go/src/github.com/dangarmol/go-notes/03-goroutines-examples/main.go:129 +0xa4
==================

Found 3 data race(s)
```
