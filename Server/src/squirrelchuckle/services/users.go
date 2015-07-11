package services
import (
	"sync"
	
	"squirrelchuckle/core"
	"gopkg.in/mgo.v2"
	"errors"
)

type UserService struct {
	alive bool

	Users map[string]User `email to user`

	*mgo.Collection
	sync.Mutex
}

type Avatar struct {

}

type Subscribe struct {

}

type UserInfo struct {
	Department string	`json:"department" bson:"department"`
	Avatar
}

type User struct {
	Email 		string			`json:"email" bson:"_id"`
	Name 		string  		`json:"name" bson:"name"`
	UserInfo
	Devices 	[]*DeviceToken

	core.PrivilegeLevel			`json:"privilege" bson:"privilege"`
	dirty bool
}

func (this *UserService) Alive() bool {
	return this.alive
}

func (this *UserService) Depends() []string {
	return nil
}

func (this *UserService) Name() string {
	return "UserService"
}

func (this *UserService) Initialize() error {
	this.Lock()
	defer this.Unlock()
	if this.alive {
		return nil
	}

	step 	:= core.DbQueryLimit
	cursor 	:= step
	users 	:= make([]User, 0, step)
	this.Users = make(map[string]User)

	this.Collection = core.SquirrelApp.DB("squirrel").C("user")
	q := this.Find(nil)

	q.Limit(cursor).All(&users)
	for ; len(users) != 0; cursor += step {
		for _, user := range users {
			this.Users[user.Email] = user
		}
		q.Skip(cursor).Limit(step).All(&users)
	}
	this.alive = true
	return nil
}

func (this *UserService) UnInitialize() {
	this.Lock()
	defer this.Unlock()
	if !this.alive {
		return
	}

	this.Collection = nil
	this.alive = false
}

func (this *UserService) validate(user *User) bool {
	if _, ok := this.Users[user.Email]; ok {
		return false
	}

	return true
}

func (this *UserService) BulkAdd(users []*User) error {
	length := len(users)
	switch length {
	case 1:
		return this.Add(users[0])
	case 0:
		return nil
	}

	var validCount int
	valid := make([]*User, length)
	for _, user := range users {
		if this.validate(user) {
			valid[validCount] = makeUser(user)
			validCount += 1
		}
	}

	valid = valid[:validCount]
	if validCount != len(users) {
		core.SquirrelApp.Error("Some users not added %v", length - validCount)
	}

	bulk := this.Bulk()
	bulk.Unordered()
	bulk.Insert(valid)
	_, err := bulk.Run()

	return err
}

func makeUser(user *User) *User {
	user.PrivilegeLevel = core.Normal
	if user.Devices == nil {
		user.Devices = make([]*DeviceToken, 1)
	}
	return user
}

func (this *UserService) Add(user *User) error {
	if !this.validate(user) {
		return errors.New("invalid user")
	}
	
	return this.Insert(makeUser(user))
}