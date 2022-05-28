package server

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	addr := ":8080"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is testing"))
	})
	err := Server(addr)
	t.Fatal(err)
}
