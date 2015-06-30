package services
import (
	"sync"
	"strconv"

	"squirrelchuckle/database"
	"squirrelchuckle/settings"
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
	push_unreachable_count int
	conn_reachable_tick int
	user_id DbKeyType
}

type DeviceTokenService struct {
	alive bool
	deviceTokens map[int32]*DeviceToken
	changedTokens map[int32]*DeviceToken `dirty tokens, false for stall entries, true for new entries`
	sync.Mutex
	currentTick int
	staleTick int

	staleTokens []int32
}


var defaultPushUnreachableTolerance, defaultConnUnreachableTolerance int = 20, 7 * 24
var pushUnreachableTolerance, connUnreachableTolerance int

func (this *DeviceTokenService) Initialize() (error) {
	this.Lock()
	defer this.Unlock()

	if this.alive {
		return nil
	}

	this.changedTokens = make(map[int32]*DeviceToken)
	this.staleTokens = make([]int32, 10000)

	var err error

	pushUnreachableTolerance, err = strconv.Atoi(settings.AppConfig["push_unreachable_tolerance"])
	if err != nil {
		pushUnreachableTolerance = defaultPushUnreachableTolerance
	}

	connUnreachableTolerance, err = strconv.Atoi(settings.AppConfig["conn_unreachable_tolerance"])
	if err != nil {
		connUnreachableTolerance = defaultConnUnreachableTolerance
	}

	c := database.MSession.DB("squirrel").C("device_token")
	q := c.Find(nil)
	iterator := q.Iter()

	var result []struct { token int32 }
	iterator.All(&result)
	iterator.Close()

	this.deviceTokens = make(map[int32]*DeviceToken, len(result))

	for _, token := range result {
		this.deviceTokens[token.token] = new(DeviceToken)
	}

	this.alive = true
	return nil
}

func (this *DeviceTokenService) UnInitialize() {
	this.Lock()
	defer this.Unlock()

	if !this.alive {
		return
	}

	this.flush()
	this.deviceTokens = nil
	this.alive = false
}

func (this *DeviceTokenService) flush() {
	this.Lock()
	defer this.Unlock()
}

func (this *DeviceTokenService) Alive() bool {
	return this.alive
}

/****************************************************
	methods
 ****************************************************/

func (this *DeviceTokenService) Add(tokens []int32) {
	for _, token := range tokens {
		e := this.deviceTokens[token]
		if e != nil {
			e.conn_reachable_tick = this.currentTick
			e.push_unreachable_count = 0
		} else {
			this.deviceTokens[token] = new (DeviceToken)
		}
	}
}

func (this *DeviceTokenService) Touch (tokens []int32) {
	for _, token := range tokens {
		e := this.deviceTokens[token]
		if e != nil {
			e.push_unreachable_count = 0
		}
	}
}

func (this *DeviceTokenService) Disconnect (tokens [] int32) {
	for _, token := range tokens {
		e := this.deviceTokens[token]
		if e != nil {
			e.conn_reachable_tick = this.currentTick
		}
	}
}

// TODO: tick overflow
func (this *DeviceTokenService) tick() {
	this.currentTick++
	this.staleTick = this.currentTick - connUnreachableTolerance
}

func (this *DeviceTokenService) Stale() {

	// TODO: tick overflow
	for k, v := range this.deviceTokens {
		if v.push_unreachable_count >= pushUnreachableTolerance && v.conn_reachable_tick >= this.staleTick {
			this.staleTokens = append(this.staleTokens, k)
		}
	}

	for _, v := range this.staleTokens {
		delete(this.deviceTokens, v)
	}

	this.staleTokens = this.staleTokens[0:0]
}