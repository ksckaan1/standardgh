package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ksckaan1/standardgh"
)

func main() {
	http.HandleFunc("POST /standard/{id}", standardgh.GH(standardHandler))
	http.HandleFunc("GET /stream", standardgh.GHforSSE(5*time.Second, chatHandler))

	fmt.Println("listening")
	fmt.Println("Try: curl -N 'http://localhost:3030/stream?room=general'")
	err := http.ListenAndServe(":3030", nil)
	if err != nil {
		panic(err)
	}
}
