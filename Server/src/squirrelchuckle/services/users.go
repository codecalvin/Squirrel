package services
import (
	"sync"
	"errors"

	"gopkg.in/mgo.v2"
	"encoding/hex"

	"github.com/twinj/uuid"
	"squirrelchuckle/core"
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

type ClassBriefItem struct {
	ElementType_UniqueKey string
	ElementType_ClassName string
	ElementType_ClassTime string
}

type ClassItem struct {
	ElementType_UniqueKey        string
	ElementType_ClassName        string
	ElementType_ClassTime        string
	ElementType_ClassTeacher     string
	ElementType_ClassStudent     string
	ElementType_ClassDescription string
	RegisterUsers                map[string]string
}

type User struct {
	AdsName     string          `json:"login_name" bson:"_id"`
	Email 		string			`json:"email" bson:"email"`
	Name 		string  		`json:"name" bson:"name"`
	AdsPass     string        	`json:"-" bson:"-"`
	Password  	string 			`json:"-" bson:"password"`
	UserInfo
	Devices 	[]*DeviceToken
	Classes		map[string]*ClassBriefItem

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

	if !core.SquirrelApp.Auth(&user.AdsName, &user.AdsPass) {
		return errors.New("invalid password")
	}
	user.Password = uuid.NewV4().String()

	this.Users[user.AdsName] = makeUser(user)
	this.newUsers = append(this.newUsers, user)

	return nil
}

func findDevice(user *User, device *DeviceToken) (int, *DeviceToken) {
	for i, d := range(user.Devices) {
		if device == d {
			return i, d
		}
	}
	return -1, nil
}
//
//func (this *UserService) TransferDevice(user *User, deviceToken string) {
//	if device, ok := deviceTokenService.deviceTokens[deviceToken]; ok {
//		if device.UserID == user.AdsName {
//			// touch return
//			touchDevice(device)
//		} else {
//			// transfer device
//		}
//
//	} else {
//		this.AddDevice()
//	}
//}


func (this *UserService) AddWithDevice(user *User, deviceToken string) (*User, error) {
	if !this.validate(user) {
		return nil, errors.New("invalid user")
	}

	if !core.SquirrelApp.Auth(&user.AdsName, &user.AdsPass) {
		return nil, errors.New("invalid password")
	}
	user.Password = uuid.NewV4().String()

	this.Lock()
	defer this.Unlock()

	if content, err := hex.DecodeString(deviceToken); err != nil || len(content) != APNS_TOKEN_SIZE {
		return nil, errors.New("Error device token")
	}

	this.Users[user.AdsName] = makeUser(user)
	if deviceTokenService != nil {
		if err := deviceTokenService.Add(user.AdsName, deviceToken); err != nil {
			user.Devices = nil
			return user, err
		}
	}

	return user, nil
}

// d is retrieved from DeviceTokenService
func (this *UserService) RemoveDevice(d *DeviceToken) {
	if user, ok := this.Users[d.UserID]; ok {
		for i := 0 ; i < len(user.Devices); i++ {
			if user.Devices[i] == d {
				// as d is new from DeviceTokenService
				// so we only need to comparing the reference
				// || user.Devices[i].DeviceTokenId == d.DeviceTokenId {
				user.Devices = append(user.Devices[:i], user.Devices[i+1:]...)
				user.dirty = true
				break
			}
		}
	}
}

func (this *UserService) Auth(adsName, adsPass *string) bool{
	return core.SquirrelApp.Auth(adsName, adsPass)
}

// d is retrieved from DeviceTokenService
func (this *UserService) AddDevice(d *DeviceToken) {
	if user, ok := this.Users[d.UserID]; ok {
		if user.Devices == nil {
			user.Devices = []*DeviceToken {d}
		} else {
			user.Devices = append(user.Devices, d)
		}
		user.dirty = true
	}
}
