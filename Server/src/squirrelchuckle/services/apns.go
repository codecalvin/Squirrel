package services

import (
	"crypto/tls"

	"time"
	"sync"
	"fmt"
	"strconv"
)

var apnConn *APNSConnection
var feedConn *tls.Conn

type ApplePushService struct {
	alive 		bool

	ticker 		*time.Ticker
	stopChan 	chan bool
	sync.Mutex
}

var applePushService *ApplePushService
func (this *ApplePushService) Alive() bool {
	return this.alive
}

func (this *ApplePushService) Depends() []string {
	return nil
}

func (this *ApplePushService) Name() string {
	return "ApplePushService"
}

func (this *ApplePushService) Initialize() error {
	this.Lock()
	defer this.Unlock()

	if this.alive {
		return nil
	}
	var err error
	if cert, err := tls.LoadX509KeyPair("conf/SquirrelCert.pem.private", "conf/SquirrelKey.u.pem.private"); err == nil {
		var tlsConn *tls.Conn
		tlsConf := &tls.Config { Certificates: []tls.Certificate{cert} }
		if tlsConn, err = tls.Dial("tcp", "gateway.sandbox.push.apple.com:2195", tlsConf); err != nil {
			goto stale
		}
		if apnConn, err = socketAPNSConnection(tlsConn, &APNSConfig{}); err != nil {
			tlsConn.Close()
			goto stale
		}

		if feedConn, err = tls.Dial("tcp", "feedback.sandbox.push.apple.com:2196", tlsConf); err != nil {
			goto stale
		}

		//success:
		this.alive = true
		applePushService = this
		return nil
	}

	this.stopChan = make(chan bool)

	// ticker
	this.ticker = time.NewTicker(time.Second * 10)
	go func () {
	loop:
		for {
			select {
			case <-	this.ticker.C:
				apnsTicker(this)
			case <- this.stopChan:
				break loop
			}
		}
	} ()

stale:
	return err
}

func apnsTicker(this *ApplePushService) {
	if elements, err := readFromFeedbackService(feedConn); err != nil {
		for element := elements.Front(); element != nil; element = element.Next() {
			value, _ := element.Value.(*FeedbackResponse)
			fmt.Println(value.Token)
		}
	}
}

func (this *ApplePushService) UnInitialize() {
	this.Lock()
	defer this.Unlock()

	if !this.alive {
		return
	}

	this.stopChan <- true
	this.ticker.Stop()
	apnConn.Disconnect()
	apnConn = nil
	feedConn.Close()
	feedConn = nil
	this.alive = false
	applePushService = nil
}

func (this *ApplePushService) NotifySimple(deviceToken, message, badge string) {
	var nBadge *Badge
	if v, err := strconv.ParseInt(badge, 10, 32); err == nil {
		nBadge = NewBadge(int(v))
	}
	payload := &Payload {
		AlertText	: 	message,
		Badge		:	nBadge,
		Sound		: 	"default",
		Category	: 	"Information",
		Token		: 	deviceToken,
		ContentAvailable: 1,
	}
	apnConn.SendChannel <- payload
}