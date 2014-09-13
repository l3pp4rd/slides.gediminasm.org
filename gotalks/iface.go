package main

import "fmt"

type secretCar struct {
	Name string
}

func (c *secretCar) String() string {
	return fmt.Sprintf("Secret car named: %s", c.Name)
}

func New(name string) fmt.Stringer {
	return &secretCar{name}
}

func main() {
	printables := make([]fmt.Stringer, 0)
	printables = append(printables, New("BMW"), New("Audi"))
	for _, printable := range printables {
		fmt.Printf("A stringer '%s'\n", printable)
	}
}
