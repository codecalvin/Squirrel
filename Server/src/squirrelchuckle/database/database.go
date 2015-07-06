package database

import (
	"gopkg.in/mgo.v2"

	"squirrelchuckle/settings"
)

var MSession *mgo.Session

var DataBaseNameString string
var DataBaseClassCollectionNameString string
var DataBaseUserCollectionNameString string

func init() {
	var err error
	MSession, err = mgo.Dial(settings.AppDatabaseAddr)
	
	if err != nil {
		panic(err)
	}
	
	DataBaseNameString = "squirrel" 
	DataBaseClassCollectionNameString = "class" 
	DataBaseUserCollectionNameString = "user"
}

func Close() {
	if MSession != nil {
		MSession.Close()
	}
	MSession = nil
}