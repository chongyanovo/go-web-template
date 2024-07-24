package internal

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config 配置文件
type Config struct {
	Server *Server `mapstructure:"server" json:"server" yaml:"server"`
	Mysql  *Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Zap    *Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
}

// NewConfig 读取配置文件
func NewConfig(v *viper.Viper) *Config {
	config := Config{}
	if err := v.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %v", err))
	}
	return &config
}
