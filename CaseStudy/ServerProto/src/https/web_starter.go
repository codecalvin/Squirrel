package main

import (
	"net/http"
	"fmt"
	"flag"
)

var httpsEnable = flag.Bool("https.enabled", false, "enable https service")

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func main()  {

	t := NewUnstartedServer(http.HandlerFunc(defaultHandler))

	flag.Parse()

	if *httpsEnable {
		t.StartTLS()
	} else {
		t.Start()
	}

	defer t.Close()
}