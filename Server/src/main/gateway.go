package main

import (
	"net/http"
	"fmt"
	"sync"

	"github.com/gorilla/mux"

	"server/v1"
	"server"
	"utils"
)

var ss server.SquirrelServer
var gl sync.Mutex

// default handler, only for present
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Squirrel!")
	fmt.Fprintln(w, "Current version %s", ss.Version())
	fmt.Fprintln(w, "valid versions %v", utils.UInt8Keys(ss.ORs))
}

// TODO: json configuration file for current version APIs
func initServer() {
	gl.Lock()
	defer gl.Unlock()

	if ss.Status != server.Dead {
		return
	}

	ss.ORs = map[string]*mux.Router{}
	ss.AddRouter(mux.NewRouter(), 1, 1)
}

func unloadServer() {
	gl.Lock()
	defer gl.Unlock()
}

func main() {
	initServer()

	ss.R.Handle("/", HomeHandler)
	v1.InstallHandlers(ss.SquirrelRouter)
}
