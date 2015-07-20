package services

import (
	"crypto/tls"

	"time"
	"sync"
	"fmt"
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
	fmt.Println("apns ticker...")
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

var deviceToken string = "ed8bf85a32fd6a3e6339531ae6fe90824d67df100a11f5dfe7f22ae0f92a50b8"
//						 "ed8bf85a32fd6a3e6339531ae6fe90824d67df100a11f5dfe7f22ae0f92a50b8"
func (this *ApplePushService) TestAPN() error {
	//write token
	payload := &Payload {
		AlertText: "A new class coming",
		Badge: NewBadge(0),
		ContentAvailable: 1,
		Sound: "default",
		Token: deviceToken,
	}
	fmt.Println(payload.Badge)
	apnConn.SendChannel <- payload
	return nil
}