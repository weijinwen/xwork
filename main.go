package main

import (
	"xwork/BootStrap/Config"
	"xwork/BootStrap/DbBase"
	"xwork/BootStrap/Server"
)

func main() {
	Config.InitConfig()
	DbBase.InitDb()
	//DbBase.InitCache()
	Server.InitIris()
}
