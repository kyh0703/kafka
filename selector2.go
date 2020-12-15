package main

import (
	"fmt"
	"time"
)

func process(ch chan string) {
	time.Sleep(10 * time.Second)
	ch <- "process successful"
}

func scheduling() {

}

func main() {
	ch := make(chan string)
	go process(ch)
	for {
		time.Sleep(1 * time.Second)
		select {
		case v := <-ch:
			fmt.Println("received value: ", v)
			return
		default:
			fmt.Println("no Value received")
		}

		scheduling()
	}
}
