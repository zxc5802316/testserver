package main

import (
	"fmt"
	"net/http"
	"week3/server"
)

func main() {
	addr := ":8080"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	err := server.Server(addr)
	fmt.Println(err)
}
