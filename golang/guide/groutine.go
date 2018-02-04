package main

import (
	"time"
	"fmt"
)

func Say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}
func main() {
	go Say("world")
	Say("Hello")
}
