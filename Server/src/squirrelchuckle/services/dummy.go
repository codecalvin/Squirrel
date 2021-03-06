package services

import "sync"

type DummyService struct {
	alive bool
	sync.Mutex
}

var dummyService *DummyService
func (this *DummyService) Alive() bool {
	return this.alive
}

func (this *DummyService) Depends() []string {
	return nil
}

func (this *DummyService) Name() string {
	return "DummyService"
}

func (this *DummyService) Initialize() error {
	this.Lock()
	defer this.Unlock()

	if this.alive {
		return nil
	}

	this.alive = true
	dummyService = this
	return nil
}

func (this *DummyService) UnInitialize() {
	this.Lock()
	defer this.Unlock()

	if !this.alive {
		return
	}

	this.alive = false
	dummyService = nil
}
