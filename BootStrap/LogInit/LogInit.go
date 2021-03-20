package LogInit

import (
	"io"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//@param app iris.Application
//初始化日志系统
func LogInit(app *iris.Application) *iris.Application {
	//系统日志配置
	serverLogFile := newServerLogFile()
	app.Logger().SetLevel(viper.GetString("LogLevel"))
	app.Logger().SetOutput(io.MultiWriter(serverLogFile, os.Stdout))

	//日志工具配置
	webLogFile := newWebLofFile()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(webLogFile)
	level, _ := logrus.ParseLevel(viper.GetString("LogLevel"))
	logrus.SetLevel(level)

	return app
}

//系统日志文件
func newServerLogFile() *os.File {
	filename := viper.GetString("LogFile") + "server.log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

//日志工具文件
func newWebLofFile() *os.File {
	filename := viper.GetString("LogFile") + "web.log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}
