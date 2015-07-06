package core
import "sync"

type ServiceManager struct {
	services map[string]ServiceInterface
	spawnChan chan ServiceInterface
	termChan chan bool
	postInitHandler map[string] []PostInitFunc
	sync.Mutex
	alive bool
}

func spawnRouting(this *ServiceManager) {
loop:
	for {
		select {
		case spawn := <-this.spawnChan:
			name := spawn.Name()
			if err := spawn.Initialize(); err == nil {
				this.Lock()
				this.services[name] = spawn
				this.Unlock()
			} else {
				SquirrelApp.Critical("[ServiceManager]: Failed to initialize service %v. Error: %v", name, err)
			}
		case <- this.termChan:
			break loop
		}
	}
}

func (this *ServiceManager) Alive() bool {
	return true
}

func (this *ServiceManager) Name() string {
	return "ServiceManager"
}

func (this *ServiceManager) Initialize() error {
	this.services = make(map[string]ServiceInterface)
	this.spawnChan = make(chan ServiceInterface)
	this.termChan = make(chan bool)
	
	this.postInitHandler = make(map[string] []PostInitFunc)
	this.alive = true

	go spawnRouting(this)

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
		initHandlers = make([]PostInitFunc, 5)
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
	
	switch updateServices() {
	case 2:
		this.spawnChan <- service
		fallthrough
	case 1:
		SquirrelApp.Info("[ServiceManager]: Register service :%v", name)
		return true
	default:
		return false
	}
}