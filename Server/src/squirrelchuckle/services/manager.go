package services

type ServiceManager struct {
	Services map[string]ServiceInterface
	postInitHandler map[string] []PostInitFunc
}

var mgrInstance *ServiceManager

func init() {
	mgrInstance = new(ServiceManager)
	mgrInstance.Services = make(map[string]ServiceInterface)
	mgrInstance.postInitHandler = make(map[string] []PostInitFunc)
}

func GetManager() *ServiceManager {
	return mgrInstance
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
	mgrInstance.Services["APNService"] = &APNService{}
	mgrInstance.Services["DeviceTokenService"] = &DeviceTokenService{}
	mgrInstance.Services["UserService"] = &UserService{}

	for k := range mgrInstance.Services {
		mgrInstance.Services[k].Initialize()
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
	for k := range mgrInstance.Services {
		mgrInstance.Services[k].UnInitialize()
	}
}

func (this *ServiceManager) GetServiceByName (name string) ServiceInterface {
	return mgrInstance.Services[name]
}