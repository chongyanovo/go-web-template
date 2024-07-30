package bootstrap

import (
	"github.com/ChongYanOvO/go-web-template/bootstrap/internal"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	App Application
)

type Application struct {
	Config *internal.Config
	Viper  *viper.Viper
	Mysql  *gorm.DB
	Logger *zap.Logger
}

// InitApp 初始化 App
func InitApp() Application {
	App.Viper = internal.NewViper()
	App.Config = internal.NewConfig(App.Viper)
	App.Logger = App.Config.Zap.NewZap()
	zap.ReplaceGlobals(App.Logger)
	App.Mysql = internal.NewMysql(App.Config.Mysql)
	return App
}
