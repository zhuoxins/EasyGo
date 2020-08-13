package cache

import (
	"EasyGo/cache/drivers"
	"EasyGo/config"
	"EasyGo/tools/exception"
	"errors"
)

var creatorList map[string]creator

type creator func(conf config.Result) (Provider, error)

type Cache interface {
	Put(key string, value interface{}, expire ...int) bool
	Get(key string) interface{}
	Forget(key string) bool
	Remember(key string, expire int, callValue func() interface{}) interface{}
	Has(key string) bool
	Err() *exception.Exception
}

type Provider interface {
	SetValue(key string, value interface{}, expire ...int) (bool, error)
	GetValue(key string) (interface{}, error)
	Remove(key string) bool
	HasKey(key string) bool
}

//实例化redis驱动
func createRedisDriver(conf config.Result) (Provider, error) {
	if conf == nil {
		return nil, errors.New("invalid redis config")
	}
	redis, err := drivers.NewRedis(conf)
	if err != nil {
		return nil, err
	}
	return redis, nil
}
