package core
import (
	"sync"
	"fmt"
)

type ServiceManager struct {
	services 		map[string] ServiceInterface
	spawnChan 		chan ServiceInterface
	matureChan 		chan bool
	termChan 		chan bool

	postInitHandler map[string] []PostInitFunc
	// pending service name => service
	pendingService 	map[string] ServiceInterface
	// pending service name => depends service names
	pendingChain    []ServiceInterface
	sync.Mutex
	alive bool
}

func findMature(manager *ServiceManager) ServiceInterface {
	manager.Lock()
	defer manager.Unlock()
	var si ServiceInterface

	// find from the last one
	for i := len(manager.pendingChain) - 1; i >= 0; i-- {
		si = manager.pendingChain[i]
		for _, name := range si.Depends() {
			if _, ok := manager.services[name]; !ok {
				si = nil
				break
			}
		}

		if si != nil {
			manager.pendingChain = append(manager.pendingChain[:i], manager.pendingChain[i+1:]...)
			break
		}
	}
	return si
}

func matureRouting(this *ServiceManager) {
loop:
	for {
		select {
		case <- this.matureChan:
			// check service ready to initialize
			var si ServiceInterface
			if si = findMature(this); si == nil {
				continue
			}

			name := si.Name()
			if err := si.Initialize(); err == nil {
				fmt.Println("Initialize service ", name)
				SquirrelApp.Info("Initialize service %v", name)
				this.Lock()
				this.services[name] = si
				this.Unlock()
				go func() {
					this.matureChan <- true
				}()
			} else {
				SquirrelApp.Error("[ServiceManager]: Failed to initialize service %v. Error: %v", name, err)
				continue
			}
		case <- this.termChan:
			break loop
		}
	}
}

func (this *ServiceManager) Alive() bool {
	return true
}

func (this *ServiceManager) Depends() [] string {
	return nil
}

func (this *ServiceManager) Name() string {
	return "ServiceManager"
}

func (this *ServiceManager) Initialize() error {
	this.termChan = make(chan bool)
	this.matureChan = make(chan bool)
	this.spawnChan = make(chan ServiceInterface)
	this.services = make(map[string]ServiceInterface)

	this.pendingService = make(map[string]ServiceInterface)
	this.pendingChain = make([]ServiceInterface, 0, 5)
	this.postInitHandler = make(map[string] []PostInitFunc)
	this.alive = true

	go matureRouting(this)

	return nil
}

func (this *ServiceManager) UnInitialize() {
	close(this.termChan)
	close(this.spawnChan)
	this.Lock()
	defer this.Unlock()

	for k := range this.services {
		SquirrelApp.Info("Unloading service %v", this.services[k].Name())
		this.services[k].UnInitialize()
	}
	this.services = nil
	this.alive = false
}

func (this *ServiceManager) AddPostInitHandler(service string, function PostInitFunc) {
	initHandlers := this.postInitHandler[service]
	if initHandlers == nil {
		initHandlers = make([]PostInitFunc, 0, 5)
		this.postInitHandler[service] = initHandlers
	}

	initHandlers = append(initHandlers, function)
}

func (this *ServiceManager) GetServiceByName (name string) ServiceInterface {
	return this.services[name]
}

func (this *ServiceManager) RegisterService (service ServiceInterface) bool {
	if !this.alive {
		return false
	}
	name := service.Name()

	updateServices := func () int {
		this.Lock()
		defer this.Unlock()
		if _, ok := this.services[name]; ok {
			SquirrelApp.Warning("[ServiceManager]: Already register service :%v", name)
			return 0
		} else if service.Alive() {
			this.services[name] = service
			return 1
		}
		return 2
	}
	
	code := updateServices()
	switch code {
	case 2:
		this.Lock()
		this.pendingChain = append(this.pendingChain, service)
		this.Unlock()
		fallthrough
	case 1:
		this.matureChan <- true
		SquirrelApp.Info("[ServiceManager]: Register service :%v", name)
		return true
	default:
		return false
	}
}