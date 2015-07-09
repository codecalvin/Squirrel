package services
import (
	"sync"
	"strconv"

	"time"
	"squirrelchuckle/core"
	"gopkg.in/mgo.v2"
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

	core.DeviceID						`json:"device_id" bson:"device_id"`
	UserID 				core.DbDefKey 	`json:"user_id" bson:"user_id"`
	DeviceTokenId 		[8]byte 		`json:"id" bson:"_id"`
}

type DeviceTokenService struct {
	deviceTokens map[uint32]*DeviceToken
	changedTokens map[uint32]*DeviceToken `dirty tokens, false for stall entries, true for new entries`

	*mgo.Collection
	currentTick,staleTick int
	staleTokens []uint32
	alive bool

	uniqueId uint

	ticker *time.Ticker
	stopChan chan bool
	sync.Mutex
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

	this.changedTokens = make(map[uint32]*DeviceToken)
	this.staleTokens = make([]uint32, 0, 10000)
	this.stopChan = make(chan bool)
	
	var err error
	if pushUnreachableTolerance, err = strconv.Atoi(core.SquirrelApp.RunConfig("push_unreachable_tolerance")); err != nil {
		pushUnreachableTolerance = defaultPushUnreachableTolerance
	}

	if connUnreachableTolerance, err = strconv.Atoi(core.SquirrelApp.RunConfig("conn_unreachable_tolerance")); err != nil {
		connUnreachableTolerance = defaultConnUnreachableTolerance
	}

	// setup unique id
	this.uniqueId = core.SquirrelApp.UniqueId("device_token")
	this.Collection = core.SquirrelApp.DB("squirrel").C("device_token")
	q := this.Find(nil)
	iterator := q.Iter()

	var result []struct { token uint32 }
	iterator.All(&result)
	iterator.Close()

	this.deviceTokens = make(map[uint32]*DeviceToken, len(result))

	for _, token := range result {
		this.deviceTokens[token.token] = new(DeviceToken)
	}

	// ticker
	this.ticker = time.NewTicker(time.Millisecond * 100)
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

func (this *DeviceTokenService) getUniqueId(inc uint) uint {
	this.Lock()
	defer this.Unlock()
	curId := this.uniqueId
	this.uniqueId += inc
	return curId
}

func (this *DeviceTokenService) Add(token *DeviceToken) error {
	token.connReachableTick = this.currentTick
	token.DeviceID = core.DeviceID(this.getUniqueId(1))
	return this.Insert(token)
}

func (this *DeviceTokenService) BulkAdd(tokens []*DeviceToken) error {
	length := len(tokens)

	switch length {
	case 1:
		return this.Add(tokens[0])
	case 0:
		return nil
	}

	curId := this.getUniqueId(uint(length))
	bulk := this.Bulk()
	bulk.Unordered()

	// update token
	for _, token := range tokens {
		token.connReachableTick = this.currentTick
		token.DeviceID = core.DeviceID(curId)
		curId += 1
	}
	bulk.Insert(tokens)

	var err error
	if _, err = bulk.Run(); err != nil {
		core.SquirrelApp.Error(err.Error())
	}
	return err
}

func (this *DeviceTokenService) Touch (tokens []uint32) {
	for _, token := range tokens {
		if e := this.deviceTokens[token]; e != nil {
			e.pushUnreachableTick = 0
		}
	}
}

func (this *DeviceTokenService) Disconnect (tokens []uint32) {
	for _, token := range tokens {
		if e := this.deviceTokens[token]; e != nil {
			e.connReachableTick = this.currentTick
		}
	}
}

func (this *DeviceTokenService) Stale() {
	// TODO: tick overflow
	for k, v := range this.deviceTokens {
		if v.pushUnreachableTick >= pushUnreachableTolerance && v.connReachableTick >= this.staleTick {
			this.staleTokens = append(this.staleTokens, k)
		}
	}

	for _, v := range this.staleTokens {
		delete(this.deviceTokens, v)
	}

	this.staleTokens = this.staleTokens[0:0]
}

func (this *DeviceTokenService) TestAPN() {
	if _, ok := core.SquirrelApp.GetServiceByName("ApplePushService").(*ApplePushService); ok {
	}
}