package core
import (
	"github.com/astaxie/beego/logs"
)

type Squirrel struct {
	*AppSetting
	*Database
	*ServiceManager
	auth *AuthService
	*logs.BeeLogger
}

var (
	SquirrelApp *Squirrel
)

var CloseChan chan bool
func init() {
	CloseChan = make(chan bool)
	SquirrelApp = NewApp()
	if err := SquirrelApp.Initialize(); err != nil {
		SquirrelApp.Fatal("SquirrelApp initialize failed. Fatal: %v", err)
	}
}

func Run() {
	go func() {
		defer CloseApp()
		<-CloseChan
	} ()
}

func NewApp() *Squirrel {
	return &Squirrel{}
}

func CloseApp() {
	if SquirrelApp != nil {
		SquirrelApp.UnInitialize()
		SquirrelApp = nil
	}
}

func (this *Squirrel) Fatal(format string, v... interface{}) {
	this.Critical(format, v...)
	panic("Fatal")
}

func (this *Squirrel) Initialize() error {
	this.BeeLogger = logs.NewLogger(10000)
	this.AppSetting = &AppSetting{}
	this.Database = &Database{}
	this.auth = &AuthService{}

	this.ServiceManager = &ServiceManager{}
	this.ServiceManager.Initialize()
	this.ServiceManager.RegisterService(this.Database)
	this.ServiceManager.RegisterService(this.AppSetting)
	this.ServiceManager.RegisterService(this.auth)

	return nil
}

func (this *Squirrel) UnInitialize() {
	this.ServiceManager.UnInitialize()
	this.Database.UnInitialize()
	this.AppSetting.UnInitialize()
	this.BeeLogger.Close()
}

func (this *Squirrel) Auth(name, password *string) bool {
	return this.auth.Auth(name, password)
}