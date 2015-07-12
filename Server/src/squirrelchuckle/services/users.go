package services
import (
	"sync"
	
	"squirrelchuckle/core"
	"gopkg.in/mgo.v2"
	"errors"
)

type UserService struct {
	Users 			map[string]*User `email to user`

	updateUsers 	[]*User
	newUsers		[]*User

	*mgo.Collection
	sync.Mutex
	alive 			bool
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
	AdsName     string          `json:"login_name" bson:"_id"`
	Email 		string			`json:"email" bson:"email"`
	Name 		string  		`json:"name" bson:"name"`
	Password  	string 			`json:"-" bson:"password"`
	UserInfo
	Devices 	[]*DeviceToken

	core.PrivilegeLevel			`json:"privilege" bson:"privilege"`
	dirty		bool
}

func (this *UserService) Alive() bool {
	return this.alive
}

func (this *UserService) Depends() []string {
	return []string { "AuthService", "Database" }
}

func (this *UserService) Name() string {
	return "UserService"
}

var userService *UserService
func (this *UserService) Initialize() error {
	this.Lock()
	defer this.Unlock()
	if this.alive {
		return nil
	}

	step 	:= core.DbQueryLimit
	cursor 	:= step
	users 	:= make([]*User, 0, step)
	this.Users = make(map[string]*User)

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
	userService = this
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
	userService = nil
}

func (this *UserService) validate(user *User) bool {
	if _, ok := this.Users[user.AdsName]; ok {
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

	var err error
	if validCount != len(users) {
		core.SquirrelApp.Error("Some users not added %v", length - validCount)
		err = errors.New("Some users not added")
	}

	valid = valid[:validCount]
	this.Lock()
	defer this.Unlock()

	this.newUsers = append(this.newUsers, valid...)
	for _, user := range valid {
		this.Users[user.AdsName] = user
	}

	return err
}

func makeUser(user *User) *User {
	user.PrivilegeLevel = core.Normal
//	if user.Devices == nil {
//		user.Devices = make([]*DeviceToken, 1)
//	}
	return user
}

func (this *UserService) Add(user *User) error {
	if !this.validate(user) {
		return errors.New("invalid user")
	}

	this.Lock()
	defer this.Unlock()
	this.Users[user.AdsName] = makeUser(user)
	this.newUsers = append(this.newUsers, user)

	return nil
}

func (this *UserService) AddWithDevice(user *User, deviceToken string) error {
	if !this.validate(user) {
		return errors.New("invalid user")
	}

	this.Lock()
	defer this.Unlock()
	this.Users[user.AdsName] = makeUser(user)
	user.Devices = make([]*DeviceToken, 1)

	newDevice		:= makeNewDevice(user.AdsName, deviceToken)
	this.newUsers 	= append(this.newUsers, user)
	user.Devices[0] = newDevice
	if deviceTokenService != nil {
		deviceTokenService.deviceTokens[deviceToken] = newDevice
	}

	return nil
}

func (this *UserService) RemoveDevice(d *DeviceToken) {
	if user, ok := this.Users[d.UserID]; ok {
		i := 0
		for ; i < len(user.Devices) && user.Devices[i] != d ; i++ {}
		user.Devices = append(user.Devices[:i], user.Devices[i+1:]...)
		user.dirty = true
	}
}

func (this *UserService) AddDevice(d *DeviceToken) {
	if user, ok := this.Users[d.UserID]; ok {
		i := 0
		for ; i < len(user.Devices) && user.Devices[i] != d ; i++ {}
		user.Devices = append(user.Devices[:i], user.Devices[i+1:]...)
		user.dirty = true
	}
}
