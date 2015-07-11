package services

import (
	"crypto/tls"

	"fmt"
	"encoding/json"
	"encoding/binary"
	"bytes"
)


type simpleAps struct {
	Alert            string
	Badge            int
	Sound            string
	Category         string
	ContentAvailable int
}

// APNS item ids
const (
	DEVICE_TOKEN_ID = 1
	PAYLOAD_ID = 2
	NOTIFICATION_CENTER_ID = 3
	EXPIRATION_DATE_ID = 4
	PRIORITY_ID = 5
)

const (
	DEVICE_TOKEN_LEN = 32

	TEST_ID_NUM = 0xfeeeee
)

func (s *simpleAps) MarshalJSON() ([]byte, error) {
	toMarshal := make(map[string]interface{})
	msgMarshal := make(map[string]interface{})
	if s.Alert != "" {
		toMarshal["alert"] = s.Alert
	}

	toMarshal["badge"] = s.Badge
	if s.Sound != "" {
		toMarshal["sound"] = s.Sound
	}
	if s.Category != "" {
		toMarshal["category"] = s.Category
	}
	if s.ContentAvailable != 0 {
		toMarshal["content-available"] = s.ContentAvailable
	}

	msgMarshal["aps"] = toMarshal
	return json.Marshal(msgMarshal)
}


var APNConn *tls.Conn

//ed8bf85a 32fd6a3e 6339531a e6fe9082 4d67df10 0a11f5df e7f22ae0 f92a50b8
//var testDeviceToken string = "ed8bf85a32fd6a3e6339531ae6fe90824d67df100a11f5dfe7f22ae0f92a50b8"

type ApplePushService struct {
	alive bool
}

func (this *ApplePushService) Alive() bool {
	return this.alive
}

func (this *ApplePushService) Depends() []string {
	return [] string { "ApplePushService" }
}

func (this *ApplePushService) Name() string {
	return "ApplePushService"
}

func (this *ApplePushService) Initialize() error {
	if this.alive {
		return nil
	}

	var err error
	var conf *tls.Config
	if cert, err := tls.LoadX509KeyPair("conf/SquirrelCert.pem.private", "conf/SquirrelKey.u.pem.private"); err == nil {
		conf = &tls.Config { Certificates: []tls.Certificate{cert} }
		if APNConn, err = tls.Dial("tcp", "gateway.sandbox.push.apple.com:2195", conf); err != nil {
			goto stale
		}

//success:
		this.alive = true
		return nil
	}

stale:
	return err
}

func (this *ApplePushService) UnInitialize() {
	if !this.alive {
		return
	}

	APNConn.Close()
	this.alive = false
}

func (this *ApplePushService) TestAPN(tokens []DeviceToken) error {
	payload := simpleAps {
		Badge:1,
		Alert:"hello my apple",
		Sound:"default",
	}

	payloadBytes, err := payload.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Error marshalling payload %v\n", err)
	}

	if err == nil {
		return nil
	}

	//write token

	frameByteBuffer := new(bytes.Buffer)
	for _, token := range tokens {
		itemByteBuffer := new(bytes.Buffer)
		binary.Write(itemByteBuffer, binary.BigEndian, uint8(1))
		binary.Write(itemByteBuffer, binary.BigEndian, uint16(DEVICE_TOKEN_LEN))
		binary.Write(itemByteBuffer, binary.BigEndian, token)

		testBuffer := new(bytes.Buffer)
		binary.Write(testBuffer, binary.BigEndian, token)

		//write payload
		binary.Write(itemByteBuffer, binary.BigEndian, uint8(2))
		binary.Write(itemByteBuffer, binary.BigEndian, uint16(len(payloadBytes)))
		binary.Write(itemByteBuffer, binary.BigEndian, payloadBytes)
	
		//write id
		binary.Write(itemByteBuffer, binary.BigEndian, uint8(3))
		binary.Write(itemByteBuffer, binary.BigEndian, uint16(4))
		binary.Write(itemByteBuffer, binary.BigEndian, TEST_ID_NUM)
	
		//write header info and item info
		binary.Write(frameByteBuffer, binary.BigEndian, uint8(2))
		binary.Write(frameByteBuffer, binary.BigEndian, uint32(itemByteBuffer.Len()))
		itemByteBuffer.WriteTo(frameByteBuffer)
	}
	
	//write to socket
	_, writeErr := APNConn.Write(frameByteBuffer.Bytes())
	if writeErr != nil {
		fmt.Printf("Error while writing to socket \n", writeErr)
	}
	
	return writeErr
}