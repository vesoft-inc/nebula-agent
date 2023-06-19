package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/conf"
)

var C Config

type Config struct {
	Agent      string
	PluginPath string
}

// parse yaml config file
func InitConfig(configFilePath string) {
	C = Config{}
	// if has config file, load config from config file
	if _, err := os.Stat(configFilePath); err != nil {
		logrus.Warnf("config file %s not exist, use default config", configFilePath)
		return
	}
	conf.MustLoad(configFilePath, &C, conf.UseEnv())
}
