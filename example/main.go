package main

import (
	"fmt"
	"net/http"

	"github.com/ksckaan1/standardgh"
)

func main() {
	http.HandleFunc("POST /standard/{id}", standardgh.GH(standardHandler))

	fmt.Println("listening")
	err := http.ListenAndServe(":3030", nil)
	if err != nil {
		panic(err)
	}
}
