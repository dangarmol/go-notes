# Channels

Most programming languages were designed before multithreading was commonplace, and rely on external libraries and packages to somewhat enable concurrency. Go took concurrency into account since its design, thus `channels` are a tool to pass data around goroutines safely whilst preventing issues such as race conditions and other memory sharing problems.

## Introduction

- Channels are declared with the keyword `chan`.
- A common case for goroutines with channels is, if you have data that is asynchronously processed. For example, if data is generated quickly but it takes a while to process, or viceversa, if data takes a while to be generated and there are multiple generators, but then it's quickly processed.

## Channels creation

- Channels are strongly typed, so they will only allow one type of data for the messages in them.
- They need to be created using make: `make(chan int)`.

```go
var wg = sync.WaitGroup{}

func simpleChannelDemo() {
   ch := make(chan int)
   wg.Add(2)
   go func() {
      i := <-ch
      fmt.Println(i)
      wg.Done()
   }()
   go func() {
      j := 42
      ch <- j
      j = 13 // This doesn't affect the data passed to the channel.
      wg.Done()
   }()
   wg.Wait()
}
```

## Sender-Receiver deadlock

- A potential problem in these sorts of scenarios is when all goroutines are blocked.
- This may happen because the channel has only sent one element but expects to receive multiple ones.
- Channels synchronise threads by waiting for data to arrive or waiting to send data before continuing execution.
- The following example is dangerous, because if either one of the goroutines is executed one less time than in the example, a `deadlock` would happen.

```go
var wg = sync.WaitGroup{}

func potentialDeadlockDemo() {
   fmt.Println("Potential deadlock:")
   ch := make(chan int)
   for j := 0; j < 5; j++ {
      wg.Add(1)
      go func() {
         i := <-ch
         fmt.Println(i)
         wg.Done()
      }()
   }
   for j := 0; j < 5; j++ {
      wg.Add(1)
      go func() {
         ch <- 42
         wg.Done()
      }()
   }
   wg.Wait()
}
```

- `ch <- 42` pauses the execution of the goroutine until there is a space available in the channel (for unbuffered channels).
- `i := <-ch` pauses execution of the goroutine until there is data to consume in the channel.
- By default, we are working with `unbuffered channels`, which means that only one message can be in the channel at a time.
- **This behaviour can cause deadlocks both if there are more elements sent than consumed, or if there are more elements expected than received.**

## Sender-Receiver by direction

In the following example, these operations will happen:

- `i := <-ch` in "Up" is waiting, there is nothing in the channel for it.
- `ch <- 42` in "Down" puts 42 in the channel and waits for it to be consumed.
- `i := <-ch` in "Up" receives 42.
- `fmt.Println("Down:", <-ch)` in "Down" is waiting.
- `fmt.Println("Up:", i)` in "Up" `i` is printed.
- `ch <- 27` in "Up" puts 27 in the channel and waits for it to be consumed.
- `fmt.Println("Down:", <-ch)` in "Down" receives and prints 27. "Down" is `Done()`.
- "Up" is `Done()`, because 27 has been consumed.

```go
func senderAndReceiverDemo() {
   ch := make(chan int)
   wg.Add(2)
   go func() { // Up
      i := <-ch
      fmt.Println("Up:", i)
      ch <- 27
      wg.Done()
   }()
   go func() { // Down
      ch <- 42
      fmt.Println("Down:", <-ch)
      wg.Done()
   }()
   wg.Wait()
}
```

## Send-only or Receive-only channels

- It is possible for goroutines to both send and receive data, but usually it is best to constrain each goroutine to only either sending or receiving data.
- This can be achieved when declaring the goroutine with the following syntax:
  - `go func receiveStuff(ch <-chan int) {...}`: **Receive-only** channel. Will only take data from the channel.
  - `go func sendStuff(ch chan<- int) {...}`: **Send-only** channel. Will only send data into the channel.
- By using this syntax, bidirectional channels are "casted" into unidirectional channels. This is unique to channels and cannot be done with any other types.

```go
func sendReceiveOnlyDemo() {
   ch := make(chan int)
   wg.Add(2)
   go func(ch <-chan int) { // Receive only
      i := <-ch
      fmt.Println("Up:", i)
      // ch <- 27 // This would cause an error
      wg.Done()
   }(ch)
   go func(ch chan<- int) { // Send only
      ch <- 42
      // fmt.Println("Down:", <-ch) // This would cause an error
      wg.Done()
   }(ch)
   wg.Wait()
}
```

## Buffered channel

- Buffered channels can be created by specifying a buffer size on the `make(chan int, 50)` call.
- **One issue with buffered channels is that elements may be left over and not consumed and this will not cause a panic**.
- In this case, the line `ch <- 27` would cause a panic on unbuffered channels, because it would be stuck waiting to be consumed.

```go
func bufferedChannelDemo() {
   fmt.Println("Buffered channel demo:")
   ch := make(chan int, 50) // Up to 50 buffered elements
   wg.Add(2)
   go func(ch <-chan int) {
      i := <-ch
      fmt.Println(i)
      wg.Done()
   }(ch)
   go func(ch chan<- int) {
      ch <- 42
      ch <- 27
      wg.Done()
   }(ch)
   wg.Wait()
}
```

## For-range loops with Channels

- To avoid leaving data hanging without being consumed in a buffered channel, we can loop over channels to consume all data in them without knowing how many elements are there beforehand.
- The syntax for looping over channels is slightly different to the one when looping over slices, since channel messages do not have indexes.
- **Channels must be closed by the sender** when no more elements are being sent. Otherwise, the `for-loop` at the receiver would be stuck waiting forever and no more elements would be coming through.
- To close a channel, the built-in function `close()` shall be used as seen in the example.
- After a channel is closed, no more data shall be sent on it, otherwise, a `panic` will be triggered.

```go
func forRangeLoopChannelDemo() {
   fmt.Println("For-range loop over channel demo:")
   ch := make(chan int, 50)
   wg.Add(2)
   go func(ch <-chan int) {
      for i := range ch { // Loops once per element
         fmt.Println(i)
      }
      wg.Done()
   }(ch)
   go func(ch chan<- int) {
      ch <- 42
      ch <- 27
      ch <- 14
      close(ch) // Without this, the for-loop would be in a deadlock, waiting for elements forever
      wg.Done()
   }(ch)
   wg.Wait()
}
```

The `for-range` loop in channels is just syntactic sugar:

```go
for i := range ch { // Loops once per element
   fmt.Println(i)
}
// equals to
for { // Loops forever
   if i, ok := <-ch; ok {
      fmt.Println(i)
   } else {
      break // Breaks loop if channel is closed
   }
}
```

## Select statements with Channels

- **Signal-only channels** send no data through, except for the fact that an element has been sent. They can be declared with `make(chan struct{})`. Empty structs are special in Go in that they require no memory allocation.
- Signal-only channels can be used in combination with a `select` statement to improve the `for-range` functionality.
- Multiple channels can be checked for messages by using a `select` statement:

```go
var betterLogCh = make(chan logEntry, 50)
var doneCh = make(chan struct{}) // Signal-only channel, no data transmitted.

func betterLogger() {
   for {
      select {
      case entry := <-betterLogCh:
         fmt.Printf("%v - [%v] %v\n", entry.time.Format("2006-01-02T15:04:05"), entry.severity, entry.message)
      case <-doneCh:
         break
      }
   }
}

func betterLoggerDemo() {
   fmt.Println("Better way to handle signals:")
   go betterLogger()
   logCh <- logEntry{time.Now(), logInfo, "App is starting"}
   time.Sleep(2 * time.Second)
   logCh <- logEntry{time.Now(), logInfo, "App is shutting down"}
   time.Sleep(100 * time.Millisecond)
   doneCh <- struct{}{}
}
```

- It is important to note that the `select` statement will block forever until a message comes through either of the channels. If this is not desired behaviour, there can be a `default` case that does whatever is required instead.
