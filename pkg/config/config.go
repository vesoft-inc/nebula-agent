package config

import "github.com/zeromicro/go-zero/core/conf"

var C Config

type Config struct {
	HeartBeatHosts    []string
	HeartBeatInterval int64
	Agent             string
}

// parse yaml config file
func InitConfig(configFilePath string) {
	C = Config{}
	conf.MustLoad(configFilePath, &C, conf.UseEnv())
}
