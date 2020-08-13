package cache

import "EasyGo/config"

type CacheConfig struct {
	Driver string
	Prefix string
	Config config.Result
}

func (conf *CacheConfig) GetDriver() string {
	return conf.Driver
}

func (conf *CacheConfig) GetPrefix() string {
	return conf.Prefix
}

func (conf *CacheConfig) GetConfig() config.Result {
	return conf.Config
}
