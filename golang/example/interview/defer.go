package main

import "fmt"


/*
	defer是在return之前执行
	C调用return过程如下：
		1. 先给返回值赋值
		2. 调用return指令
	Golang调用return过程如下：
		1. 先给返回值赋值
		2. 调用defer指令
		3. 调用return指令
	【PS】Golang使用栈做返回值，C使用寄存器做返回值
*/
func main() {
	defer_call()
}

func defer_call() {
	defer func() {
		fmt.Printf("打印前\n")
	}()
	defer func() {
		fmt.Printf("打印中\n")
	}()
	defer func() {
		fmt.Printf("打印后\n")
	}()

	panic("触发异常")
}

//0
func f() (result int) {
	defer func() {
		result++ //执行return之前，做了一次++操作
	}()
	return 0
}

//5
func f1() (r int) {
	t := 5
	defer func() {
		t = t + 5 //这个地方不会修改r
	}()
	return t
}

//1
func f2() (r int) {
	defer func(r int) { //r是传递进去的，不会改变返回的r值
		r = r + 5
	}(r)
	return 1
}