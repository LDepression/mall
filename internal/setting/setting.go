package setting

import (
	"fmt"
	"github.com/spf13/viper"
	"mall/internal/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
	//刚才设置的环境变量 想要生效 我们必须得重启goland
}

//这里是用viper将配置文件读取到
func init() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configPrefix := "config"
	var configName string
	if debug {
		configName = fmt.Sprintf("%s_debug.yaml", configPrefix)
	} else {
		configName = fmt.Sprintf("%s_pro.yaml", configPrefix)
	}
	v := viper.New()
	v.SetConfigFile(fmt.Sprintf("config/%s", configName))
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(&global.Setting)
	if err != nil {
		panic(err)
	}
}
