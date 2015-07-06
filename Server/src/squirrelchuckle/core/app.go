package core
import "github.com/astaxie/beego/logs"

type Squirrel struct {
	*AppSetting
	*Database
	*ServiceManager
	*logs.BeeLogger
}

var (
	SquirrelApp *Squirrel
)

var CloseChan chan bool
func init() {
	CloseChan = make(chan bool)
	SquirrelApp = NewApp()
	if err := SquirrelApp.Initialize(); err == nil {
		SquirrelApp.Fatal("SquirrelApp initialize failed. Fatal: %v", err)
	}
}

func Run() {
	defer CloseApp()
	<-CloseChan
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
	this.BeeLogger = logs.NewLogger(10000).Async()
	this.AppSetting = &AppSetting{}
	this.Database = &Database{}
	this.ServiceManager = &ServiceManager{}

	err := this.AppSetting.Initialize()
	if err == nil {
		err = this.Database.Initialize()
	}
	if err == nil {
		err = this.ServiceManager.Initialize()
	}

	return err
}

func (this *Squirrel) UnInitialize() {
	this.ServiceManager.UnInitialize()
	this.Database.UnInitialize()
	this.AppSetting.UnInitialize()
	this.BeeLogger.Close()
}