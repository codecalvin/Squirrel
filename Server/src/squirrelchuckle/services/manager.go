package services

type ServiceManager struct {
	Services map[string] ServiceStuber
}

var instance *ServiceManager

func init() {
	instance = new(ServiceManager)
	instance.Services = make(map[string]ServiceStuber)
}

func GetManager() *ServiceManager {
	return instance
}

func (this *ServiceManager) Initialize() {
	instance.Services["APN_Service"] = &APNService{}
	for k := range instance.Services {
		instance.Services[k].Initialize()
	}
}

func (this *ServiceManager) UnInitialize() {
	for k := range instance.Services {
		instance.Services[k].UnInitialize()
	}
}

func (this *ServiceManager) GetServiceByName (name string) ServiceStuber {
	return instance.Services[name]
}