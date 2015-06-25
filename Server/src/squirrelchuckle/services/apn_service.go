package services

import (

	//"net"
	"crypto/tls"
)

var APNConn *tls.Conn

type APN_Service struct {
}

func (APN_Service) Initialize() (error) {
	var err error
	var conf = tls.Config {

	}
	APNConn, err = tls.Dial("tcp", "gateway.sandbox.push.apple.com:2195", &conf)
	if err != nil {
		panic("cannot reach APNs")
		return err
	}

	return nil
}

func UnInitialize() {
	APNConn.Close()
}