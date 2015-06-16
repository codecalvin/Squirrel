package database

import (
	"gopkg.in/mgo.v2"

	"squirrelchuckle/settings"
)

var MSession *mgo.Session

func init() {
	var err error
	MSession, err = mgo.Dial(settings.AppDatabaseAddr)
	
	if err != nil {
		panic(err)
	}
}

func Close() {
	if MSession != nil {
		MSession.Close()
	}
	MSession = nil
}