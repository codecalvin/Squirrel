package core

type ServiceInterface interface {
	Alive() bool
	Name() string
	Depends() []string
	Initialize() error
	UnInitialize()
}


const DbQueryLimit = 10000

type DeviceID uint
type DeviceKind uint8
type PrivilegeLevel uint16

type PostInitFunc func ()

const (
	INVALID DeviceKind = iota
	IOS = 10 + iota
	ANDROID
)

const (
	Delete PrivilegeLevel = iota
	ANONYMOUS
	Normal
	Admin
	Super
)
