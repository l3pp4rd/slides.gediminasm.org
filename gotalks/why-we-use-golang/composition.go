package main

import (
	"fmt"
	"talks/car"
)

type BmwX5 struct {
	car.Car // composed of Car
	Engine  string
}

func New() *BmwX5 {
	return &BmwX5{
		car.Car{Vendor: "BMW", Model: "X5", Doors: 5}, "V8",
	}
}

func main() {
	fmt.Printf("Car: %s\n", New().Drive(5))
}
