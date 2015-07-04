package services
import (
	"sync"
	
	"squirrelchuckle/core"
)

type UserService struct {
	alive bool

	users map[string]UserInfo `email to user`

	sync.Mutex
}

type Avatar struct {

}

type Device struct {
	core.DeviceID
	UserID core.DbKeyType
}

type Subscribe struct {

}

type UserInfo struct {
	Name string
	NickName string
	Email string
	Avatar
	Devices []Device
	Subscribes []Subscribe

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

	this.users = make(map[string]UserInfo)

	c := core.SquirrelApp.MSession.DB("squirrel").C("user")
	q := c.Find(nil)

	userInfo := make([]UserInfo, core.DbQueryLimit)

	for {
		q.Limit(core.DbQueryLimit).All(&userInfo)
		if len(userInfo) == 0 {
			break
		}
		for _, user := range userInfo {
			this.users[user.Email] = user
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

	this.alive = false
}
