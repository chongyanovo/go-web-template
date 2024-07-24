package internal

import (
	"fmt"
	"github.com/spf13/viper"
)

// NewViper 创建viper实例
func NewViper() *viper.Viper {
	v := viper.New()
	v.AddConfigPath("./config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %s \n", err))
	}
	return v
}
