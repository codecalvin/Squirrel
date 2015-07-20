package core

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

type Database struct {
	*mgo.Session
	counterCol *mgo.Collection
	sync.Mutex

	alive bool
}

type Counter struct {
	Field string	`bson:"field"`
	Count uint		`bson:"count"`
}

func (this *Database) Alive() bool {
	return this.alive
}

func (this *Database) Depends() []string {
	return []string {"AppSetting"}
}

func (this *Database) Name() string {
	return "Database"
}

func (this *Database) Initialize() error {
	if this.alive {
		return nil
	}
	
	var err error
	if this.Session, err = mgo.Dial(SquirrelApp.AppDatabaseAddr); err == nil {
		this.counterCol = this.DB("squirrel").C("counter")
		this.alive = true
	} else {
		SquirrelApp.Critical("Database Initialize failed. Error: %v", err)
	}

	return err
}

func (this *Database) UnInitialize() {
	if this.alive {
		this.Close()
	}

	this.alive = false
}

func (this *Database) UniqueId(field string) uint {
	change := mgo.Change{
		Update: bson.M{"$inc": bson.M{"count" : 1} },
		ReturnNew: true,
		Upsert: true,
	}

	var result Counter
	_, err := this.counterCol.Find(bson.M{ "field": field }).Apply(change, &result)
	if err != nil {
		SquirrelApp.Error(err.Error())
	}
	return result.Count
}
