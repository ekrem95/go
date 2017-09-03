package main

import (
	"fmt"
	"time"
)

func getMessages(msg string, delay time.Duration) <-chan string {
	c := make(chan string)

	go func() {
		for i := 1; i <= 3; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Millisecond * delay)
		}
	}()
	return c
}

func main() {
	c1 := getMessages("first", 300)
	c2 := getMessages("second", 150)
	c3 := getMessages("third", 10)

	// for i := 1; i <= 3; i++ {
	// 	println(<-c1)
	// 	println(<-c2)
	// 	println(<-c3)
	// }

	for i := 1; i <= 9; i++ {
		select {
		case msg := <-c1:
			println(msg)
		case msg := <-c2:
			println(msg)
		case msg := <-c3:
			println(msg)
		}
	}
}
