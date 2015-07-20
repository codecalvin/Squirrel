package services

import (
	"crypto/tls"

	"fmt"
	"encoding/json"
	"encoding/binary"
	"bytes"
)

var apnConn *APNSConnection
type ApplePushService struct {
	alive bool
}

var applePushService *ApplePushService
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
	if cert, err := tls.LoadX509KeyPair("conf/SquirrelCert.pem.private", "conf/SquirrelKey.u.pem.private"); err == nil {
		var tlsConn *tls.Conn
		tlsConn, err = tls.Dial("tcp", "gateway.sandbox.push.apple.com:2195", &tls.Config { Certificates: []tls.Certificate{cert} })
		if err != nil {
			goto stale
		}

		if apnConn, err = socketAPNSConnection(tlsConn, &APNSConfig{}); err != nil {
			goto stale
		}
//success:
		this.alive = true
		applePushService = this
		return nil
	}

stale:
	return err
}

func (this *ApplePushService) UnInitialize() {
	if !this.alive {
		return
	}

	apnConn.Disconnect()
	apnConn = nil
	this.alive = false
	applePushService = nil
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
	apnConn.SendChannel
	frameByteBuffer := new(bytes.Buffer)
	for _, token := range tokens {
		itemByteBuffer := new(bytes.Buffer)
		binary.Write(itemByteBuffer, binary.BigEndian, uint8(1))
		binary.Write(itemByteBuffer, binary.BigEndian, uint16(DEVICE_TOKEN_LEN))
		binary.Write(itemByteBuffer, binary.BigEndian, token)

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