package services

type ServiceManager struct {
	Services map[string]ServiceInterface
	postInitHandler map[string] []PostInitFunc
}

var instance *ServiceManager

func init() {
	instance = new(ServiceManager)
	instance.Services = make(map[string]ServiceInterface)
	instance.postInitHandler = make(map[string] []PostInitFunc)
}

func GetManager() *ServiceManager {
	return instance
}

func (this *ServiceManager) AddPostInitHandler(service string, function PostInitFunc) {
	initHandlers := this.postInitHandler[service]
	if initHandlers == nil {
		initHandlers = make([]PostInitFunc, 5)
		this.postInitHandler[service] = initHandlers
	}

	initHandlers = append(initHandlers, function)
}

func (this *ServiceManager) Initialize() {
	instance.Services["APNService"] = &APNService{}
	instance.Services["DeviceTokenService"] = &DeviceTokenService{}
	instance.Services["UserService"] = &UserService{}

	for k := range instance.Services {
		instance.Services[k].Initialize()
		initHandlers := this.postInitHandler[k]
		if initHandlers == nil {
			for _, v := range initHandlers {
				v()
			}
		}
	}
	this.postInitHandler = nil
}

func (this *ServiceManager) UnInitialize() {
	for k := range instance.Services {
		instance.Services[k].UnInitialize()
	}
}

func (this *ServiceManager) GetServiceByName (name string) ServiceInterface {
	return instance.Services[name]
}