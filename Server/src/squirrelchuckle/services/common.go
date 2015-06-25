package services

type ServiceStub interface {
	Initialize() (error)
	UnInitialize()
}
