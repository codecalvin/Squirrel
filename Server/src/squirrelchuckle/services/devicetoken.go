package services
import (
	"sync"
	"strconv"

	"time"
	"squirrelchuckle/core"
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
	pushUnreachableCount int
	connReachableTick int
	userId core.DbKeyType
}

type DeviceTokenService struct {
	alive bool
	deviceTokens map[int32]*DeviceToken
	changedTokens map[int32]*DeviceToken `dirty tokens, false for stall entries, true for new entries`
	sync.Mutex
	currentTick int
	staleTick int

	staleTokens []int32
	ticker *time.Ticker
	stopChan chan bool
}


var defaultPushUnreachableTolerance, defaultConnUnreachableTolerance int = 20, 7 * 24
var pushUnreachableTolerance, connUnreachableTolerance int


func (this *DeviceTokenService) Alive() bool {
	return this.alive
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

	this.changedTokens = make(map[int32]*DeviceToken)
	this.staleTokens = make([]int32, 10000)
	this.stopChan = make(chan bool)
	
	var err error
	if pushUnreachableTolerance, err = strconv.Atoi(core.SquirrelApp.RunConfig("push_unreachable_tolerance")); err != nil {
		pushUnreachableTolerance = defaultPushUnreachableTolerance
	}

	if connUnreachableTolerance, err = strconv.Atoi(core.SquirrelApp.RunConfig("conn_unreachable_tolerance")); err != nil {
		connUnreachableTolerance = defaultConnUnreachableTolerance
	}

	c := core.SquirrelApp.MSession.DB("squirrel").C("device_token")
	q := c.Find(nil)
	iterator := q.Iter()

	var result []struct { token int32 }
	iterator.All(&result)
	iterator.Close()

	this.deviceTokens = make(map[int32]*DeviceToken, len(result))

	for _, token := range result {
		this.deviceTokens[token.token] = new(DeviceToken)
	}

	// ticker
	this.ticker = time.NewTicker(time.Second)
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
	return nil
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
}

func ticker(this *DeviceTokenService) {
	this.currentTick++
	this.staleTick = this.currentTick - connUnreachableTolerance
	this.Stale()
}

func (this *DeviceTokenService) flush() {
	this.Lock()
	defer this.Unlock()
}

/****************************************************
	methods
 ****************************************************/

func (this *DeviceTokenService) Add(tokens []int32) {
	for _, token := range tokens {
		e := this.deviceTokens[token]
		if e != nil {
			e.connReachableTick = this.currentTick
			e.pushUnreachableCount = 0
		} else {
			this.deviceTokens[token] = new (DeviceToken)
		}
	}
}

func (this *DeviceTokenService) Touch (tokens []int32) {
	for _, token := range tokens {
		e := this.deviceTokens[token]
		if e != nil {
			e.pushUnreachableCount = 0
		}
	}
}

func (this *DeviceTokenService) Disconnect (tokens [] int32) {
	for _, token := range tokens {
		e := this.deviceTokens[token]
		if e != nil {
			e.connReachableTick = this.currentTick
		}
	}
}

func (this *DeviceTokenService) Stale() {

	// TODO: tick overflow
	for k, v := range this.deviceTokens {
		if v.pushUnreachableCount >= pushUnreachableTolerance && v.connReachableTick >= this.staleTick {
			this.staleTokens = append(this.staleTokens, k)
		}
	}

	for _, v := range this.staleTokens {
		delete(this.deviceTokens, v)
	}

	this.staleTokens = this.staleTokens[0:0]
}