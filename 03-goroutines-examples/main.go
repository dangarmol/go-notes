package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}
var counter = 0

var mtx = sync.RWMutex{}
var counterMutex = 0

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

// Now all values are correctly printed in order
// However, the parallelism is ruined and the code will run sequentially
// as if there were no goroutines at all.
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

// There are still repeated values, but at least they come in order
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

func sayHello() {
	fmt.Printf("Hello #%v\n", counter)
	wg.Done()
}

func increment() {
	counter++
	wg.Done()
}

// Pretty much random values are printed in this example, they are not even sorted
func forWaitGroupExample() {
	fmt.Println("Unsynchronised 'for' WaitGroup example:")
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go sayHello()
		go increment()
	}
	wg.Wait()
}

func waitGroupExample() {
	startTime := time.Now()
	fmt.Println("WaitGroup example:")
	var msg = "Hello"
	wg.Add(2) // How many goroutines should be awaited.
	go func() {
		fmt.Println(msg) // Will print "Goodbye".
		wg.Done()
	}()
	go func(msg string) {
		fmt.Println(msg) // Will print "Hello".
		wg.Done()
	}(msg)
	msg = "Goodbye"
	wg.Wait() // Waits until the WaitGroup count is zero.
	fmt.Printf("Time waited: %v\n", time.Since(startTime))
}

func sleepExample() {
	startTime := time.Now()
	fmt.Println("Sleep example:")
	var msg = "Hello"
	go func() {
		fmt.Println(msg) // Will print "Goodbye".
	}()
	go func(msg string) {
		fmt.Println(msg) // Will print "Hello".
	}(msg)
	msg = "Goodbye"
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Time waited: %v\n", time.Since(startTime))
}

func main() {
	fmt.Println("Threads available:", runtime.GOMAXPROCS(-1))
	sleepExample()
	waitGroupExample()
	forWaitGroupExample()
	mutexExample()
	betterMutexExample()
	fmt.Println("Threads available:", runtime.GOMAXPROCS(-1))
}
