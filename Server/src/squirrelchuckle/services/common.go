package services

type ServiceStuber interface {
	Alive() bool
	Initialize() error
	UnInitialize()
}
