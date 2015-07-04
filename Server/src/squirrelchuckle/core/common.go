package core

type ServiceInterface interface {
	Alive() bool
	Name() string
	Initialize() error
	UnInitialize()
}

type DbKeyType int

const DbQueryLimit = 10000

type DeviceID uint8

type PostInitFunc func ()

const (
	NIL_DBKey DbKeyType = iota // nil db key

	IPHONE DeviceID = iota
	IPAD
	IOS = 10 + iota
	ANDROID
	WEB
)
