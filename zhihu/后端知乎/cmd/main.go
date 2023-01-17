package main

import (
	"main/boot"
	"main/utils"
)

func main() {
	boot.ViperSetup(utils.Path())
	boot.Loggersetup()
	boot.MysqlDBSetup()
	boot.RedisSetup()
	boot.ServerSetup()
}
