package utility

import (
	"github.com/astaxie/beego/utils"
)

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	return utils.FileExists(name)
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}