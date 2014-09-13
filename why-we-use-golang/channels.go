package main

import (
	"fmt"
	"time"
)

func main() {
	var one, two = make(chan string), make(chan string)
	go func() {
		for {
			time.Sleep(time.Second)
			two <- "msg1"
			one <- "msg2"
		}
	}()
	for {
		select {
		case msg := <-one:
			fmt.Printf("received %s from chan: one\n", msg)
		case msg := <-two:
			fmt.Printf("received %s from chan: two\n", msg)
		}
	}
}
