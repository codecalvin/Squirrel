package services

type ServiceStuber interface {
	Initialize() (error)
	UnInitialize()
}
