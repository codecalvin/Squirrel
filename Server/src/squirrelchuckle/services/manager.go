package services

type ServiceManager struct {
	Services map[string] *ServiceStub
}

var instance *ServiceManager

func init() {
	instance = new(ServiceManager)
	instance.Services = make(map[string]*ServiceStub)
}

func New() *ServiceManager {
	return instance
}

func (this *ServiceManager) Initialize() {
	var se := APN_Service{}
	v, _ := se.(ServiceStub)
	instance.Services["APN_Service"] = &v
}

func (this *ServiceManager) UnInitialize() {
	for k := range instance.Services {
		instance.Services[k].UnInitialize()
	}
}