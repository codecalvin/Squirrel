package services

import (
	"gopkg.in/mgo.v2"
)

type Database struct {
	MSession *mgo.Session
	alive bool
}

var dbInstance *Database

func DatabaseInstance() *Database {
	if dbInstance == nil {
		dbInstance = &Database{}
	}
	return dbInstance
}

func (this *Database) UnInitialize() {
	if this.alive {
		this.MSession.Close()
	}
	
	this.alive = false
}

func (this *Database) Initialize() error {
	if this.alive {
		return nil
	}
	
	var err error
	this.MSession, err = mgo.Dial(AppSettingInstance().AppDatabaseAddr)

	if err == nil {
		this.alive = true
	}
	
	return err
}

func (this *Database) Alive() bool {
	return this.alive
}