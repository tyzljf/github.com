package main

import "fmt"

type People struct {

}

func (p *People) showA() {
	fmt.Printf("people showA ...\n")
}

func (p *People) showB() {
	fmt.Printf("people showB ...\n")
}


type Teacher struct {
	People
}

func (t *Teacher) showB() {
	fmt.Printf("teacher showB ...\n")
}

func main() {
	t := &Teacher{}
	t.showA()
}