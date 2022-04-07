package main

import (
	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}

func manualForLoopChannelDemo() {
	fmt.Println("Manual for loop over channel demo:")
	ch := make(chan int, 50)
	wg.Add(2)
	go func(ch <-chan int) {
		for { // Loops forever
			if i, ok := <-ch; ok {
				fmt.Println(i)
			} else {
				break // Breaks loop if channel is closed
			}
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
		ch <- 27 // This would cause an error on unbuffered channels, and would not be consumed in buffered ones
		wg.Done()
	}(ch)
	wg.Wait()
}

func sendReceiveOnlyDemo() {
	fmt.Println("Send-only + receive-only demo:")
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

func senderAndReceiverDemo() {
	fmt.Println("Sender and receiver demo:")
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

func potentialDeadlockDemo() {
	fmt.Println("Potential deadlock:")
	ch := make(chan int)
	for j := 0; j < 5; j++ {
		wg.Add(1)
		go func() {
			i := <-ch // Pauses execution of the goroutine until there is data in the channel
			fmt.Println(i)
			wg.Done()
		}()
	}
	for j := 0; j < 5; j++ {
		wg.Add(1)
		go func() {
			ch <- 42 // Pauses execution of the goroutine until there is space in the channel
			wg.Done()
		}()
	}
	wg.Wait()
}

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
		j = 13 // This doesn't affect the data passed to the channel, will still be 42
		wg.Done()
	}()
	wg.Wait()
}

func main() {
	simpleChannelDemo()
	potentialDeadlockDemo()
	senderAndReceiverDemo()
	sendReceiveOnlyDemo()
	bufferedChannelDemo()
	forRangeLoopChannelDemo()
	manualForLoopChannelDemo()
}
