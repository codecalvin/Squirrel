package main

import (
	"net/http"
	"fmt"
)

func main() {
	ts := NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	ts.Start()
	//defer ts.Close()
}
