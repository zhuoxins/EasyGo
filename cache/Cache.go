package cache

import (
	"EasyGo/config"
	"EasyGo/kernel"
	"EasyGo/tools/exception"
	"errors"
	"fmt"
)

func init() {
	creatorList = make(map[string]creator)
	creatorList = map[string]creator{
		"redis": createRedisDriver,
	}
	cacheConf := config.Get("cache")
	if cacheConf.GetField("driver").String() != "" {
		kernel.Register("cache").Do(func() {
			conf := &CacheConfig{
				Driver: cacheConf.GetField("driver").String(),
				Prefix: cacheConf.GetField("prefix").String(),
			}
			conf.Config = config.Get("database." + conf.Driver)
			if conf.Config == nil {
				panic("cache: invalid driver config")
			}
			err := NewCache(conf)
			if err != nil {
				fmt.Println(err.Error())
			}
		})
	}
}

var handler Cache

type CacheManager struct {
	driver    string
	prefix    string
	store     map[string]Provider
	provider  Provider
	exception *exception.Exception
}

func NewCache(conf *CacheConfig) error {
	cache := &CacheManager{
		store: make(map[string]Provider),
	}
	err := cache.init(conf)
	if err != nil {
		return err
	}
	handler = cache
	return err
}

func (c *CacheManager) init(conf *CacheConfig) error {
	err := c.createDriver(conf)
	if err != nil {
		return err
	}
	err = c.setDriver(conf.GetDriver())
	if err != nil {
		return err
	}
	c.prefix = conf.GetPrefix()
	return nil
}

func (c *CacheManager) setDriver(driver string) error {
	if provider, ok := c.store[driver]; ok {
		c.driver = driver
		c.provider = provider
		return nil
	} else {
		return errors.New("cache driver not exists")
	}
}

func (c *CacheManager) createDriver(configer *CacheConfig) error {
	if creators, ok := creatorList[configer.GetDriver()]; ok {
		provider, err := creators(configer.GetConfig())
		if err != nil {
			return err
		}
		c.store[configer.GetDriver()] = provider
	}
	return errors.New("cache: invalid cache driver")
}

func (c *CacheManager) jointKey(key string) string {
	return c.prefix + key
}

//插入缓存数据
func (c *CacheManager) Put(key string, value interface{}, expire ...int) bool {
	res, err := c.provider.SetValue(c.jointKey(key), value, expire...)
	if err != nil {
		c.putError(err)
	}
	return res
}

//获取缓存数据
func (c *CacheManager) Get(key string) interface{} {
	res, err := c.provider.GetValue(c.jointKey(key))
	if err != nil {
		c.putError(err)
	}
	return res
}

//销毁缓存数据
func (c *CacheManager) Forget(key string) bool {
	return c.provider.Remove(c.jointKey(key))
}

//缓存存在读缓存,否则先获取数据再存入缓存
func (c *CacheManager) Remember(key string, expire int, callValue func() interface{}) interface{} {
	if c.provider.HasKey(key) {
		return c.Get(key)
	}
	res := callValue()
	_ = c.Put(key, res, expire)
	return res
}

func (c *CacheManager) Has(key string) bool {
	return c.provider.HasKey(c.jointKey(key))
}

//记录执行错误
func (c *CacheManager) putError(err error) {
	if c.exception == nil {
		c.exception = exception.New()
	}
	c.exception.Put(err)
}

func (c *CacheManager) Err() *exception.Exception {
	if c.exception != nil {
		err := c.exception.Clone()
		c.exception.Reset()
		return err
	}
	return exception.New()
}
