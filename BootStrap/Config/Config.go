package Config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//配置文件
func InitConfig() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./Config")

	//服务器环境配置
	viper.SetConfigName("env")
	err := viper.MergeInConfig()
	if err != nil {
		logrus.Fatalln(err)
	}
	//服务器环境配置
	env := viper.GetString("Env")

	//服务器基本配置
	viper.SetConfigName(fmt.Sprintf("server_%s", env))
	err = viper.MergeInConfig()
	if err != nil {
		logrus.Fatalln(err)
	}

	//缓存配置
	viper.SetConfigName(fmt.Sprintf("cache_%s", env))
	err = viper.MergeInConfig()
	if err != nil {
		logrus.Fatalln(err)
	}

	//数据库配置
	viper.SetConfigName(fmt.Sprintf("db_%s", env))
	err = viper.MergeInConfig()
	if err != nil {
		logrus.Fatalln(err)
	}
}
