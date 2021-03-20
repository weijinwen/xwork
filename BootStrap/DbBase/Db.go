package DbBase

import (
	"xwork/Extend/DB"
	_ "xwork/Extend/DB/Mysql"
	"os"

	"github.com/sirupsen/logrus"
)

// 初始化数据库
func InitDb() {
	//初始化主库
	orm := DB.NewOrm("default")
	if orm == nil {
		logrus.Error("db: server connection fail name default")
		os.Exit(0)
	}
}

//关闭数据库
func CloseDB() {
	//关闭主库
	DB.CloseDB("default")
}
