package DB

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DbInterface interface {
	GetOrm() *gorm.DB
	StartAndGC(server string) error
}

type Instance func() DbInterface

var adapters = make(map[string]Instance)

var servers = make(map[string]*gorm.DB)

//初始注册数据库
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("db: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("db: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

//初始化orm
func NewOrm(server string) (orm *gorm.DB) {
	var ok bool
	var err error
	config := viper.GetStringMapString(fmt.Sprintf("db.%s", server))
	adaptersName, ok := config["dialect"]
	if !ok {
		panic(fmt.Sprintf("db: server is null name %q", server))
	}
	if orm, ok = servers[server]; ok && orm != nil {
		return orm
	}
	instanceFunc, ok := adapters[adaptersName]
	if !ok {
		panic(fmt.Sprintf("db: unknown adapter name %q (forgot to import?)", config["server"]))
	}
	_, ok = config["url"]
	if !ok {
		panic(fmt.Sprintf("db: unknown url name %q", server))
	}
	adapter := instanceFunc()
	err = adapter.StartAndGC(server)
	if err != nil {
		panic(fmt.Sprintf("db: name %q err %s", server, err.Error()))
	}
	orm = adapter.GetOrm()
	servers[server] = orm
	return orm
}

// 关闭数据连接
func CloseDB(server string) {
	var orm *gorm.DB
	config := viper.GetStringMapString(fmt.Sprintf("db.%s", server))
	_, ok := config["dialect"]
	if !ok {
		logrus.Errorf("db: server is null name %q", server)
	}
	if orm, ok = servers[server]; !ok || orm == nil {
		return
	}
	if err := orm.Close(); nil != err {
		logrus.Errorf("Disconnect from database failed: %s", err.Error())
	}
	logrus.Info("Disconnect from database succeed")
}
