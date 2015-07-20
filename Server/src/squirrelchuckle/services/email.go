package services

import "sync"

type EmailService struct {
	alive bool
	sync.Mutex
}

var emailService *EmailService
func (this *EmailService) Alive() bool {
	return this.alive
}

func (this *EmailService) Depends() []string {
	return []string { "AuthService" }
}

func (this *EmailService) Name() string {
	return "EmailService"
}

func (this *EmailService) Initialize() error {
	this.Lock()
	defer this.Unlock()

	if this.alive {
		return nil
	}

	this.alive = true
	emailService = this
	return nil
}

func (this *EmailService) UnInitialize() {
	this.Lock()
	defer this.Unlock()

	if !this.alive {
		return
	}

	this.alive = false
	emailService = nil
}
