package services
import "sync"

type DummyService struct {
	alive bool
	sync.Mutex
}

func (this *DummyService) Initialize() error {
	this.Lock()
	defer this.Unlock()

	if this.alive {
		return nil
	}

	this.alive = true
	return nil
}

func (this *DummyService) UnInitialize() {
	this.Lock()
	defer this.Unlock()

	if !this.alive {
		return
	}

	this.alive = false
}

func (this *DummyService) Alive() bool {
	return this.alive
}