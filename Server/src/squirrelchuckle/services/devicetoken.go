package services
import (
	"sync"
	"strconv"

	"time"
	"squirrelchuckle/core"
	"gopkg.in/mgo.v2"
	"errors"
)

/*
	when service starts, a ticker is running to increase the global counter.
	every device token entry has two fields:
		push_unreachable_count 'puc' : used to record unreached push count since last success
		conn_reachable_tick 'crt' : used to record last global counter web-socket disconnected

	These two fields are cooperated to stale device token, disable push services

	If 'puc' greater than 'push_unreachable_tolerance', then push service disabled.
	Meanwhile, if global counter - 'crt' greater than 'conn_reachable_tick', device token is stale.

	'puc' and 'crt' are not stored into database.
 */

type DeviceToken struct {
	pushUnreachableTick int
	connReachableTick 	int

	UserID 				string	 	`json:"-" bson:"-"`
	DeviceTokenId 		string	 	`json:"id" bson:"_id"`
}

var currentTick int

type DeviceTokenService struct {
	deviceTokens 		map[string]*DeviceToken
	staleTokens         []*DeviceToken
	*mgo.Collection
	staleTick  			int

	alive bool
	ticker 				*time.Ticker
	stopChan 			chan bool
	sync.Mutex
}

var defaultPushUnreachableTolerance, defaultConnUnreachableTolerance int = 20, 7 * 24
var pushUnreachableTolerance, connUnreachableTolerance int

var deviceTokenService *DeviceTokenService

func (this *DeviceTokenService) Alive() bool {
	return this.alive
}

func (this *DeviceTokenService) Depends() []string {
	return []string { "UserService", "AppSetting" }
}

func (this *DeviceTokenService) Name() string {
	return "DeviceTokenService"
}

func (this *DeviceTokenService) Initialize() (error) {
	this.Lock()
	defer this.Unlock()

	if this.alive {
		return nil
	}

	this.stopChan = make(chan bool)
	
	var err error
	if pushUnreachableTolerance, err = strconv.Atoi(core.SquirrelApp.RunConfig("push_unreachable_tolerance")); err != nil {
		pushUnreachableTolerance = defaultPushUnreachableTolerance
	}

	if connUnreachableTolerance, err = strconv.Atoi(core.SquirrelApp.RunConfig("conn_unreachable_tolerance")); err != nil {
		connUnreachableTolerance = defaultConnUnreachableTolerance
	}

	this.deviceTokens = make(map[string]*DeviceToken, len(userService.Users))
	this.staleTokens = make([]*DeviceToken, 1000)

	for _, user := range userService.Users {
		for _, device := range user.Devices {
			device.UserID = user.AdsName
			this.deviceTokens[device.DeviceTokenId] = device
		}
	}

	// ticker
	this.ticker = time.NewTicker(time.Hour)
	go func () {
loop:
		for {
			select {
			case <-this.ticker.C:
				ticker(this)
			case <- this.stopChan:
				break loop
			}
		}
	} ()

	this.alive = true
	deviceTokenService = this
	return nil
}

func makeNewDevice(adsName, device string) *DeviceToken {
	return &DeviceToken{
		UserID: 		adsName,
		DeviceTokenId: 	device,
		connReachableTick:	currentTick,
		pushUnreachableTick:0,
	}
}

func touchDevice(device *DeviceToken) {
	device.connReachableTick = currentTick
	device.pushUnreachableTick = 0
}

func (this *DeviceTokenService) UnInitialize() {
	this.Lock()
	defer this.Unlock()

	if !this.alive {
		return
	}

	this.stopChan <- true
	this.ticker.Stop()
	this.flush()
	this.deviceTokens = nil
	this.alive = false
	deviceTokenService = nil
}

func ticker(this *DeviceTokenService) {
	currentTick++
	this.staleTick = currentTick - connUnreachableTolerance
	this.Stale()
}

func (this *DeviceTokenService) flush() {
	this.Lock()
	defer this.Unlock()
}

/****************************************************
	methods
 ****************************************************/

func (this *DeviceTokenService) Add(userName, device string) error {
	if user, ok := userService.Users[userName]; ok {
		return this.addToExistUser(user, device)
	} else {
		// not find the user
		return errors.New("invalid user name in DeviceToken add")
	}
}

func (this *DeviceTokenService) addToExistUser(user *User, device string) error {
	if v, ok := this.deviceTokens[device]; ok {
		// device already exist, update data
		if v.UserID == user.AdsName {
			// update device
			touchDevice(v)
		} else {
			userService.RemoveDevice(v)
			userService.addDevice(v)
		}
	} else {
		// device doesn't exist
		token := makeNewDevice(user.AdsName, device)

		this.Lock()
		defer this.Unlock()
		this.deviceTokens[device] = token
		userService.addDevice(token)
	}
	return nil
}

func (this *DeviceTokenService) Stale() {
	// TODO: tick overflow
	for _, v := range this.deviceTokens {
		if v.pushUnreachableTick >= pushUnreachableTolerance && v.connReachableTick >= this.staleTick {
			this.staleTokens = append(this.staleTokens, v)
		}
	}

	for _, v := range this.staleTokens {
		userService.RemoveDevice(v)
		delete(this.deviceTokens, v.DeviceTokenId)
	}

	this.staleTokens = this.staleTokens[0:0]
}

func (this *DeviceTokenService) TestAPN() {
	if applePushService != nil {
		applePushService.TestAPN()
	}
}