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
	ID 	  		core.DbDefKey   `json:"id" bson:"_id,omitempty"`
	Email 		string			`json:"email" bson:"email"`
	Name 		string  		`json:"name" bson:"name"`
	UserInfo
	Devices 	[]DeviceToken
	core.PrivilegeLevel			`json:"privilege" bson:"privilege"`

	dirty bool
}

func (this *UserService) Alive() bool {
	return this.alive
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

	this.Users = make(map[string]User)

	this.Collection = core.SquirrelApp.DB("squirrel").C("user")
	q := this.Find(nil)

	userInfo := make([]User, 0, core.DbQueryLimit)

	for {
		q.Limit(core.DbQueryLimit).All(&userInfo)
		if len(userInfo) == 0 {
			break
		}
		for _, user := range userInfo {
			this.Users[user.Email] = user
		}
		q.Skip(core.DbQueryLimit)
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

	valid := make([]*User, length)
	var i int
	for _, user := range users {
		if this.validate(user) {
			valid[i] = makeUser(user)
			i += 1
		}
	}

	valid = valid[:i]
	if i != len(users) {
		core.SquirrelApp.Error("Some users not added %v", length - i)
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
		user.Devices = make([]DeviceToken, 1)
	}
	return user
}

func (this *UserService) Add(user *User) error {
	if !this.validate(user) {
		return errors.New("invalid user")
	}
	
	return this.Insert(makeUser(user))
}