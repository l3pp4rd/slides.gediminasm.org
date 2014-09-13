package main

import "fmt"

type Car struct {
	Name string
}

func (c Car) String() string {
	return fmt.Sprintf("Car named: %s", c.Name)
}

func main() {
	fmt.Printf("Some '%s'", Car{"BMW"})
}
