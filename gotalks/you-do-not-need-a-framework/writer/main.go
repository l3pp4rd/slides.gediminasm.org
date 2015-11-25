package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

func Leveled(levels map[string]int, minLevel string) Chainable {
	lowest := levels[minLevel]
	return func(w io.Writer) io.Writer {
		return WriterFunc(func(b []byte) (int, error) {
			if start := bytes.IndexByte(b, byte('[')); start != -1 {
				if end := bytes.IndexByte(b, byte(']')); end != -1 {
					if level, ok := levels[string(b[start+1:end])]; ok {
						if level >= lowest {
							return w.Write(b)
						}
						return 0, nil
					}
				}
			}
			return w.Write(b)
		})
	}
}

func main() {
	levels := map[string]int{"DEBUG": 0, "ERROR": 1}
	log.SetOutput(New(Leveled(levels, "ERROR")).Then(os.Stdout))

	log.Println("without level message")
	log.Println("[DEBUG] debug message")
	log.Println("[ERROR] error message")
}

type WriterFunc func([]byte) (int, error)

func (w WriterFunc) Write(b []byte) (int, error) {
	return w(b)
}

type Chainable func(io.Writer) io.Writer

type Chain []Chainable

func New(ch ...Chainable) Chain {
	return Chain(ch)
}

func (c Chain) Then(w io.Writer) io.Writer {
	final := w
	for i := len(c) - 1; i >= 0; i-- {
		final = c[i](final)
	}
	return final
}
