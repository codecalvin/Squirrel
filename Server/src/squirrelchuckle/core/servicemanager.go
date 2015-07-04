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
	for spawn := range this.spawnChan {
		name := spawn.Name()
		if err := spawn.Initialize(); err != nil {
			SquirrelApp.Critical("[ServiceManager]: Failed to initialize service %v. Error: %v", name, err)
		} else {
			this.Lock()
			this.services[name] = spawn
			this.Unlock()
		}
	}
	close(this.termChan)
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
	close(this.spawnChan)
	<- this.termChan

	this.Lock()
	defer this.Unlock()

	for k := range this.services {
		this.services[k].UnInitialize()
	}

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
	this.Lock()
	defer this.Unlock()
	if !this.alive {
		return false
	}

	name := service.Name()
	if _, ok := this.services[name]; !ok {
		if service.Alive() {
			this.services[name] = service
		} else {
			this.spawnChan <- service
		}
		SquirrelApp.Info("[ServiceManager]: Register service :%v", name)
		return true
	}

	SquirrelApp.Warning("[ServiceManager]: Already register service :%v", name)
	return false
}
