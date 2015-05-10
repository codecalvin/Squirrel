package main

import (
	"net/http"
	"fmt"
	"flag"
)

var httpsEnable = flag.Bool("https.enable", false, "enable https service")

func main() {
	ts := NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client. Comeon baby let's go")
	}))

	flag.Parse()
	if *httpsEnable {
		ts.StartTLS()
	} else {
		ts.Start()
	}
	defer ts.Close()
}
