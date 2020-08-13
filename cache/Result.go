package cache

import (
	"EasyGo/tools/exception"
)

type Result string

func Get(key string) Result {
	return Result(handler.Get(key).(string))
}

func Put(key string, value interface{}, expire ...int) bool {
	return handler.Put(key, value, expire...)
}

func Forget(key string) bool {
	return handler.Forget(key)
}

func Remember(key string, expire int, callValue func() interface{}) interface{} {
	return handler.Remember(key, expire, callValue)
}

func Has(key string) bool {
	return handler.Has(key)
}

func Err() *exception.Exception {
	return handler.Err()
}
