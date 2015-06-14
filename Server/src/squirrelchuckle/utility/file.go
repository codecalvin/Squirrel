package utility

import (
	"github.com/astaxie/beego/utils"
)

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	return utils.FileExists(name)
}
