package services
import (
	"squirrelchuckle/database"
	"sync"
)


type DeviceTokenService struct {
	alive bool
	deviceTokens map[int32]bool
	changedTokens map[int32]bool
	locker chan int
}

func (this *DeviceTokenService) Initialize() (error) {
	if this.alive {
		return nil
	}

	this.locker = make(chan int)

	c := database.MSession.DB("squirrel").C("device_token")
	q := c.Find(nil)
	iterator := q.Iter()

	var result []struct { token string }
	iterator.All(&result)
	iterator.Close()

	this.deviceTokens = make(map[int32]bool, len(result))

	for _, token := range result {
		this.deviceTokens[token] = true
	}

	this.alive = true
	return nil
}

func (this *DeviceTokenService) UnInitialize() {
	if !this.alive {
		return
	}

	this.flush()
	this.deviceTokens = nil
	this.alive = false
}

func (this *DeviceTokenService) flush() {
	sync.Locker()
	this.locker <- 1
}

func (this *DeviceTokenService) Alive() bool {
	return this.alive
}