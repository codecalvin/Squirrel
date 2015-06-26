package services

import (
	"crypto/tls"
)

var APNConn *tls.Conn

type APNService struct {
}

func (APNService) Initialize() (error) {
	var err error
	//tls.LoadX509KeyPair()
	var conf = tls.Config {

	}
	APNConn, err = tls.Dial("tcp", "gateway.sandbox.push.apple.com:2195", &conf)
	if err != nil {
		panic("cannot reach APNs")
		return err
	}

	return nil
}

func (APNService) UnInitialize() {
	APNConn.Close()
}
