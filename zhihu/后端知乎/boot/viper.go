package boot

import (
	"fmt"
	"github.com/spf13/viper"
	"main/app/global"
)

func ViperSetup(ConfigPath string) {
	v := viper.New()
	v.SetConfigFile(ConfigPath) //设置配置文件路径
	v.SetConfigType("yaml")     //设置配置文件类型
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("get file failed ,err:%v", err))
	}
	if err = v.Unmarshal(&global.Config); err != nil {
		//将配置文件反序列化到config结构体
		panic(fmt.Errorf("get Unmarshal failed.err:%v", err))
	}
}
