package main

import (
	"runtime"
	"sync"
	"fmt"
)

/*
	WaitGroup
		等待所有goroutine执行完成，并且阻塞主线程的执行，直到所有goroutine执行完成
*/

func main() {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i:= 0; i < 10; i++ {
		go func() {
			fmt.Printf("waitgroup1: i=%d\n", i)
			wg.Done()
		}()
	}
	for i:= 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("waitgroup2: i=%d\n", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
