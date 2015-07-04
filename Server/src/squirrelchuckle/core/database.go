package core

import (
	"gopkg.in/mgo.v2"
)

type Database struct {
	MSession *mgo.Session
	alive bool
}

func (this *Database) Alive() bool {
	return this.alive
}

func (this *Database) Name() string {
	return "Database"
}

func (this *Database) Initialize() error {
	if this.alive {
		return nil
	}
	
	var err error
	if this.MSession, err = mgo.Dial(SquirrelApp.AppDatabaseAddr); err == nil {
		this. alive = true
	} else {
		SquirrelApp.Critical("Database Initialize failed. Error: %v", err)
	}

	return err
}

func (this *Database) UnInitialize() {
	if this.alive {
		this.MSession.Close()
	}

	this.alive = false
}
