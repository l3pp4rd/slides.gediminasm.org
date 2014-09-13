package main

import (
	"fmt"
	"math/rand"
	"time"
)

func producer(name string, publisher chan string) {
	for {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Second * time.Duration(rand.Intn(5)+2))
		publisher <- fmt.Sprintf("%s => did a job", name)
	}
}
func main() {
	res := make(chan string)
	go producer("job1", res)
	go producer("job2", res)
	go producer("job3", res)
	for result := range res {
		fmt.Printf("Got result: %s\n", result)
	}
}
