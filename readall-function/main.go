package main

import (
	"fmt"
	"io"
	"log"
	"sync"
)

// MySlowReader implements the io.Reader interface
type MySlowReader struct {
	contents string
	pos      int
	mu       sync.Mutex
}

// Read reads up to len(p) bytes into p
func (m *MySlowReader) Read(p []byte) (n int, err error) {
	if m.pos+1 <= len(m.contents) {
		n := copy(p, m.contents[m.pos:m.pos+1])
		m.pos++
		return n, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pos >= len(m.contents) {
		return 0, io.EOF
	}

	n = copy(p, m.contents[m.pos:])
	m.pos += n
	return n, nil
}

func NewMySlowReader(contents string) *MySlowReader {
	return &MySlowReader{
		contents: contents,
	}
}

func main() {

	const testString = "Hello, World!"

	mySlowReaderInstance := NewMySlowReader(testString)
	out, err := io.ReadAll(mySlowReaderInstance)
	if err != nil {
		log.Fatalf("Failed to read: %v", err)
	}

	if string(out) != testString {
		log.Fatalf("Unexpected output: got %q, want %q", out, testString)
	}

	fmt.Printf("output: %s", out)
}
