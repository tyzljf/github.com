package main

import (
	"net/http"
	"fmt"
	"log"
)

type Hello struct {
}

func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, "hello")
}

func main() {
	var h Hello
	err := http.ListenAndServe("192.168.231.134:4000", h)
	if err != nil {
		fmt.Printf("listen and serve failed !!!")
		log.Fatal(err)
	}
}
