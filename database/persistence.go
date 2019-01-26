package database

import (
	"fmt"
	"time"

	"github.com/bytom/errors"
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2/bson"

	"github.com/blockcenter/config"
)

// NewMySQLDB create the mysql master/slaves cluster
func NewMySQLDB(cfg config.MySQL, connCfg config.MySQLConnCfg) (*DB, error) {
	dsnTemplate := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local"
	masterDSN := fmt.Sprintf(dsnTemplate, cfg.Master.Username, cfg.Master.Password, cfg.Master.Host, cfg.Master.Port, cfg.Master.DbName)
	db, err := open("mysql", masterDSN)
	if err != nil {
		return nil, errors.Wrap(err, "open db cluster")
	}

	db.logMode(cfg.LogMode)
	db.SetMaxOpenConns(connCfg.MaxOpenConns)
	db.SetMaxIdleConns(connCfg.MaxIdleConns)
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping db")
	}

	return db, nil
}

type RedisDB struct {
	client         *redis.Client
	cacheSeconds   int
	longCacheHours int
}

// NewRedisDB create the redis connection
func NewRedisDB(cfg config.Redis) (*RedisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Endpoint,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "ping redis")
	}

	return &RedisDB{
		client:         client,
		cacheSeconds:   cfg.CacheSeconds,
		longCacheHours: cfg.LongCacheHours,
	}, nil
}

type CacheAble interface {
	CacheKey() string
}

func (r *RedisDB) SetWithLongTerm(obj CacheAble) error {
	expiration := time.Duration(r.longCacheHours) * time.Hour
	return r.cacheSet(obj, expiration)
}

func (r *RedisDB) SetWithShortTerm(obj CacheAble) error {
	expiration := time.Duration(r.cacheSeconds) * time.Second
	return r.cacheSet(obj, expiration)
}

// pObj is the pointer to the object unmarshaled to
func (r *RedisDB) Get(key string, pObj interface{}) error {
	result, err := r.client.Get(key).Result()
	if err != nil {
		return errors.Wrap(err, "redis get")
	}

	return bson.Unmarshal([]byte(result), pObj)
}

func (r *RedisDB) cacheSet(obj CacheAble, expiration time.Duration) error {
	b, err := bson.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "bson marshal")
	}

	return r.client.Set(obj.CacheKey(), string(b), expiration).Err()
}

func (r *RedisDB) Del(key string) *redis.IntCmd {
	return r.client.Del(key)
}
