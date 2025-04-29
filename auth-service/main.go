package main

import (
	"fmt"
	goji "goji.io"
	"net/http"

	"goji.io/pat"
)

func hello(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	fmt.Fprintf(w, "Hello, %s! This is Matcha!", name)
}

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/hello/:name"), hello)

	http.ListenAndServe("localhost:8000", mux)
}
