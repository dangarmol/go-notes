package main

import (
	"fmt"
	"time"
)

const (
	logInfo    = "INFO"
	logWarning = "WARNING"
	logError   = "ERROR"
)

type logEntry struct {
	time     time.Time
	severity string
	message  string
}

var logCh = make(chan logEntry, 50)

func logger() {
	for entry := range logCh {
		fmt.Printf("%v - [%v] %v\n", entry.time.Format("2006-01-02T15:04:05"), entry.severity, entry.message)
	}
}

func loggerDemo() {
	go logger()
	defer func() {
		time.Sleep(100 * time.Millisecond)
		close(logCh) // Closes channel 100ms after the logging is finished
	}()
	logCh <- logEntry{time.Now(), logInfo, "App is starting"}
	time.Sleep(2 * time.Second)
	logCh <- logEntry{time.Now(), logInfo, "App is shutting down"}
}

var betterLogCh = make(chan logEntry, 50)
var doneCh = make(chan struct{}) // Structs with no fields require no memory allocation

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
	doneCh <- struct{}{} // Send an empty struct on the channel to indicate it can terminate
}

func main() {
	loggerDemo()
	betterLoggerDemo()
}
