package database

import (
	"gopkg.in/mgo.v2"
)

var MSession *mgo.Session

func init() {
	var err error
	MSession, err = mgo.Dial("127.0.0.1:27017")
	
	if err != nil {
		panic(err)
	}
	
	err = MSession.Ping()

	if err != nil {
		panic(err)
	}
	print("Initialized\n")
}

func Close() {
	if MSession != nil {
		MSession.Close()
	}
	MSession = nil
}