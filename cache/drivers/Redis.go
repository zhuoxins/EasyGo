package drivers

import (
	"EasyGo/config"
	"errors"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

//redis驱动
type RedisService struct {
	client redis.Conn
	conf   config.Result
}

//实例化redis驱动
func NewRedis(conf config.Result) (*RedisService, error) {
	if conf == nil {
		return nil, errors.New("invalid redis config")
	}
	r := &RedisService{}
	err := r.connect(conf)
	if err != nil {
		return nil, err
	}
	r.conf = conf
	return r, nil
}

//连接redis
func (r *RedisService) connect(conf config.Result) error {
	if pool == nil {
		pool = &redis.Pool{
			MaxIdle:         300,
			MaxActive:       100,
			IdleTimeout:     90,
			Wait:            false,
			MaxConnLifetime: 0,
			Dial: func() (conn redis.Conn, e error) {
				return redis.Dial(conf.GetField("network").String(), conf.GetField("host").String()+":"+conf.GetField("port").String())
			},
		}
	}
	r.client = pool.Get()
	if _, err := r.client.Do("auth", conf.GetField("auth").String()); err != nil {
		_ = r.client.Close()
		return err
	}
	return nil
}

func (r *RedisService) refresh() redis.Conn {
	if err := r.client.Err(); err != nil {
		_ = r.connect(r.conf)
	}
	return r.client
}

func (r *RedisService) Set(args ...interface{}) (bool, error) {
	_, err := r.refresh().Do("SET", args...)
	var ok bool
	if err == nil {
		ok = true
	}
	return ok, err
}

func (r *RedisService) SetEx(key, value interface{}, expire int) (bool, error) {
	return redis.Bool(r.refresh().Do("SET", key, value, "EX", expire))
}

func (r *RedisService) Expire(key string, expire int) (bool, error) {
	return redis.Bool(r.refresh().Do("EXPIRE", key, expire))
}

func (r *RedisService) Get(key string) (string, error) {
	return redis.String(r.refresh().Do("GET", key))
}

func (r *RedisService) Del(key string) bool {
	ok, _ := redis.Bool(r.refresh().Do("DEL", key))
	return ok
}

func (r *RedisService) Incr(key string) (int, error) {
	return redis.Int(r.refresh().Do("INCR", key))
}

func (r *RedisService) IncrBy(key string, val interface{}) (int, error) {
	return redis.Int(r.refresh().Do("INCRBY", key, val))
}

func (r *RedisService) Decr(key string) (int, error) {
	return redis.Int(r.refresh().Do("INCR", key))
}

func (r *RedisService) DecrBy(key string, val interface{}) (int, error) {
	return redis.Int(r.refresh().Do("DECRBY", key, val))
}

func (r *RedisService) Exists(key string) bool {
	ok, _ := redis.Bool(r.refresh().Do("exists", key))
	return ok
}

func (r *RedisService) HGetAll(key string) (map[string]string, error) {
	return redis.StringMap(r.refresh().Do("HGETALL", key))
}

func (r *RedisService) HSet(key, field, value string) (bool, error) {
	return redis.Bool(r.refresh().Do("HSET", key, field, value))
}

func (r *RedisService) Smembers(key string) ([]string, error) {
	return redis.Strings(r.refresh().Do("SMEMBERS", key))
}

func (r *RedisService) Srem(key string, member interface{}) bool {
	ok, _ := redis.Bool(r.refresh().Do("SREM", key, member))
	return ok
}

func (r *RedisService) Sismember(key string, member interface{}) bool {
	ok, _ := redis.Bool(r.refresh().Do("SISMEMBER", key, member))
	return ok
}

func (r *RedisService) Ttl(key string) int {
	res, _ := redis.Int(r.refresh().Do("TTL", key))
	return res
}

func (r *RedisService) Rpop(key string) (string, error) {
	return redis.String(r.refresh().Do("RPOP", key))
}

func (r *RedisService) Close() error {
	return r.client.Close()
}

//直接执行命令操作
func (r *RedisService) Command(cmd string, content ...interface{}) (interface{}, error) {
	return r.refresh().Do(cmd, content...)
}

func (r *RedisService) SetValue(key string, value interface{}, expire ...int) (bool, error) {
	if len(expire) == 0 {
		return r.Set(key, value)
	}
	expireTime := expire[0]
	return r.SetEx(key, value, expireTime)
}

func (r *RedisService) GetValue(key string) (interface{}, error) {
	return r.Get(key)
}

func (r *RedisService) Remove(key string) bool {
	return r.Del(key)
}

func (r *RedisService) HasKey(key string) bool {
	return r.Exists(key)
}
