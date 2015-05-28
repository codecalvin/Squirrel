package server

import (
	"github.com/gorilla/mux"
	"fmt"
	"sync"
)

// refers server status
type Status int

const (
	Dead Status = iota   // uninitialized
	Alive
	Evolving             // updating
	Zombie               // tear down
)

var sl sync.Locker

type SquirrelRouter struct {
	R *mux.Router "server router"
	majorVersion uint8 "api version"
	minorVersion uint8 "internal version"
}

type SquirrelServer struct {
	*SquirrelRouter "latest router"
	Status "current status"
	version uint8 "current version"
	ORs map[uint8] *SquirrelRouter "routers map"
}


func (ss *SquirrelServer) Version () string {
	return fmt.Sprintf("v%v.%v", ss.majorVersion, ss.minorVersion)
}

func (ss *SquirrelServer) AddRouter (r *mux.Router, major uint8, minor uint8) {
	n := SquirrelRouter{ r, major, minor }

	sl.Lock()
	defer sl.Unlock()

	ss.ORs[major] = &n
	if major > ss.version {
		ss.R = r
	}
}
