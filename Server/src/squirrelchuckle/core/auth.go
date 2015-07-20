package core
import (
	"sync"
	"net/smtp"
	"errors"
	"strconv"
	"crypto/tls"
	"net"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password *string) smtp.Auth {
	return &loginAuth{*username, *password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("[loginAuth] Unkown from server")
		}
	}
	return nil, nil
}

const (
	defaultMinPool = 10
	defaultMaxPool = 50
)

type AuthService struct {
	exchangeAuth bool
	minPool int
	maxPool int

	authChan chan *smtp.Client
	*tls.Config
	alive bool
	sync.Mutex
}

func (this *AuthService) Alive() bool {
	return this.alive
}

func (this *AuthService) Depends() []string {
	return []string { "AppSetting" }
}

func (this *AuthService) Name() string {
	return "AuthService"
}

func (this *AuthService) setUpConn() *smtp.Client {
	if conn, err := net.Dial("tcp", SquirrelApp.ExchangeUrl); err == nil {
		if client, err := smtp.NewClient(conn, SquirrelApp.ExchangeHost); err != nil {
			conn.Close()
			return nil
		} else {
			return client
		}
	} else {
		return nil
	}
}

func (this *AuthService) initExchangePool() error {
	var err error

	conf := new(tls.Config)
	this.authChan = make(chan *smtp.Client, this.minPool)

	conf.Certificates = make([]tls.Certificate, 1)
	conf.Certificates[0],  err = tls.LoadX509KeyPair("conf/cert.pem", "conf/key.pem")
	if err != nil {
		SquirrelApp.Fatal("[AuthService] certificate error %v", err.Error())
	}
	conf.InsecureSkipVerify = true
	this.Config = conf
	
	for i := 0; i < this.minPool; i++ {
		if nc := this.setUpConn(); nc != nil {
			this.authChan <- nc
		} else {
			err = errors.New("[AuthService] create auth client failed")
			SquirrelApp.Fatal("%v", err)
		}
	}
	return err
}

func (this *AuthService) Initialize() error {
	this.Lock()
	defer this.Unlock()
	if this.alive {
		return nil
	}

	if this.exchangeAuth = SquirrelApp.ExchangeAuth; this.exchangeAuth {
		appConfig := SquirrelApp.AppSetting.appConfig

		this.minPool = defaultMinPool
		if value, ok := appConfig["exchange_auth_min_pool"]; ok {
			if value, err := strconv.ParseInt(value, 10, 0); err == nil {
				this.minPool = int(value)
			}
		}

		this.maxPool = defaultMaxPool
		if value, ok := appConfig["exchange_auth_max_pool"]; ok {
			if value, err := strconv.ParseInt(value, 10, 0); err == nil {
				this.maxPool = int(value)
			}
		}

		this.initExchangePool()
	}

	this.alive = true
	return nil
}

func (this *AuthService) UnInitialize() {
	this.Lock()
	defer this.Unlock()

	if !this.alive {
		return
	}
	close(this.authChan)
	this.alive = false
}

func (this *AuthService) Auth(name, password *string) bool {
	var status bool
	if this.exchangeAuth {
		select {
		case client := <- this.authChan:
			if err := client.StartTLS(this.Config); err == nil {
				status = client.Auth(LoginAuth(name, password)) == nil
				client.Close()
			}
			// return authChan
			this.authChan <- client
			return status
			// TODO timeout mechanics
		}
	}
	return status
}