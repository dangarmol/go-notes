package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	interfaceCompositionDemo()
	typeConversionDemo()
	emptyInterfaceDemo()
	caseSwitchDemo()
	valueVsPointerDemo()
}

func interfaceCompositionDemo() {
	var wc WriterCloser = NewBufferedWriterCloser()
	wc.Write([]byte("Hello, this is a test."))
	wc.Write([]byte("This will come right after."))
	wc.Write([]byte("As well as this."))
	wc.Write([]byte("There is no interruption!"))
	fmt.Println("*BEFORE BUFFER FLUSH*")
	wc.Close() // Flushes the last contents (< 8 bytes)
	fmt.Println("*AFTER BUFFER FLUSH*")
	wc.Write([]byte("Let's start over..."))
	wc.Write([]byte("YIKES!"))
	wc.Close() // Flushes the last contents (< 8 bytes)
}

func typeConversionDemo() {
	var wc WriterCloser = NewBufferedWriterCloser()
	wc.Write([]byte("Hello there, this is another test."))
	wc.Close()

	// Type conversion from `WriterCloser` to `*BufferedWriterCloser`.
	bwc := wc.(*BufferedWriterCloser)
	bwc.Write([]byte("We are now trying type conversion."))
	bwc.Write([]byte("This is now a *BufferedWriterCloser."))
	bwc.Close()

	// The type conversion shown above is risky, because it may not succeed.
	// ior := wc.(io.Reader)
	// fmt.Println(ior)
	// If run, the two lines above will cause a `panic`:
	// panic: interface conversion:
	// *main.BufferedWriterCloser is not io.Reader: missing method Read

	ior, ok := wc.(io.Reader)
	if ok {
		fmt.Println(ior)
	} else {
		fmt.Println("Conversion failed")
	}
}

func emptyInterfaceDemo() {
	var myObj interface{} = NewBufferedWriterCloser()
	if wc, ok := myObj.(WriterCloser); ok {
		wc.Write([]byte("Just one more test using type conversion "))
		wc.Write([]byte("onto an empty interface."))
		wc.Close()
	} else {
		fmt.Println("Conversion failed")
	}
}

func caseSwitchDemo() {
	var i interface{} = 0
	switch i.(type) {
	case int:
		fmt.Println("'i' is an integer")
	case string:
		fmt.Println("'i' is a string")
	default:
		fmt.Println("'i' is something else")
	}
}

func valueVsPointerDemo() {
	// Nothing is persisted
	var vwc WriterCloser = ValueWriterCloser{}
	fmt.Println("Using value as value")
	fmt.Println(vwc) // {}
	fmt.Println(vwc.Write([]byte("First")))
	fmt.Println(vwc) // {}
	fmt.Println(vwc.Write([]byte("Second")))
	fmt.Println(vwc) // {}

	// Nothing is persisted
	var vwc2 WriterCloser = &ValueWriterCloser{}
	fmt.Println("Using value as pointer")
	fmt.Println(vwc2) // &{}
	fmt.Println(vwc2.Write([]byte("First")))
	fmt.Println(vwc2) // &{}
	fmt.Println(vwc2.Write([]byte("Second")))
	fmt.Println(vwc2) // &{}

	// In this case, the internal state is persisted!
	var pwc WriterCloser = &PointerWriterCloser{}
	fmt.Println("Using pointer as pointer")
	fmt.Println(pwc) // &{}
	fmt.Println(pwc.Write([]byte("First")))
	fmt.Println(pwc) // &{First}
	fmt.Println(pwc.Write([]byte("Second")))
	fmt.Println(pwc) // &{Second}
}

type Writer interface {
	Write([]byte) (int, error)
}

type Closer interface {
	Close() error
}

// As long as you implement all the methods in the composed
// interface, then you also implement the composed interface.
type WriterCloser interface {
	Writer
	Closer
}

// The methods below do implement all methods described by the
// composed interface for `*BufferedWriterCloser`, therefore
// it implements `WriterCloser` and can be put into variables
// of this type.
type BufferedWriterCloser struct {
	buffer *bytes.Buffer
}

// Writes the contents to the console in increments of 8 bytes.
// It stores the data in the internal buffer, and then will keep
// printing until the data is less than 8 bytes.
// If the data is less than 8 bytes, it will not be consumed here,
// it will be consumed only if there is more data or if Close()
// is called.
func (bwc *BufferedWriterCloser) Write(data []byte) (int, error) {
	n, err := bwc.buffer.Write(data) // Appends data to buffer
	if err != nil {
		return 0, err
	}
	v := make([]byte, 8) // Temporary empty slice (8 bytes capacity)
	for bwc.buffer.Len() > 8 {
		_, err := bwc.buffer.Read(v)
		if err != nil {
			return 0, err
		}
		_, err = fmt.Println(string(v))
		if err != nil {
			return 0, err
		}
	}
	return n, nil
}

func (bwc *BufferedWriterCloser) Close() error {
	for bwc.buffer.Len() > 0 {
		data := bwc.buffer.Next(8)
		_, err := fmt.Println(string(data))
		if err != nil {
			return err
		}
	}
	return nil
}

// This is a constructor and is needed in order to
// initialise the internal buffer to a new buffer.
// It returns a pointer to the newly created buffer.
func NewBufferedWriterCloser() *BufferedWriterCloser {
	return &BufferedWriterCloser{
		buffer: bytes.NewBuffer([]byte{}),
	}
}

//////////////////////////////////
// Implementation using a value //
//////////////////////////////////

type ValueWriterCloser struct {
	saved string
}

func (vwc ValueWriterCloser) Write(data []byte) (int, error) {
	fmt.Println("Current saved data:", vwc.saved)
	vwc.saved = string(data)
	return len(vwc.saved), nil
}

func (vwc ValueWriterCloser) Close() error {
	return nil
}

////////////////////////////////////
// Implementation using a pointer //
////////////////////////////////////

type PointerWriterCloser struct {
	saved string
}

func (pwc *PointerWriterCloser) Write(data []byte) (int, error) {
	fmt.Println("Current saved data:", pwc.saved)
	pwc.saved = string(data)
	return len(pwc.saved), nil
}

func (pwc *PointerWriterCloser) Close() error {
	return nil
}
