package DbBase

import (
	"xwork/Extend/Cache"
	_ "xwork/Extend/Cache/Memcache"
	"os"

	"github.com/sirupsen/logrus"
)

func InitCache() {
	ca := Cache.NewCache("default")
	if ca == nil {
		logrus.Error("cache: server connection fail name default")
		os.Exit(0)
	}
}
