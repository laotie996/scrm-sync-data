package services

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"scrm-sync-data/app/config"
	"time"
)

type RedisService struct {
	context context.Context
	cancel  context.CancelFunc
	config  *config.Config //配置
	logger  *LoggerService
	error   error
	*redis.Client
	State bool
}

type Lock struct {
	key     string
	val     string
	timeout int
}

func (redisService *RedisService) Init(parentContext context.Context, config *config.Config, logger *LoggerService) {
	redisService.config = config
	redisService.logger = logger
	redisService.State = false
	redisService.context, redisService.cancel = context.WithCancel(parentContext)
	redisService.Start()
}

func (redisService *RedisService) Start() {
	fmt.Println("start redis service...", time.Now())
	redisService.logger.Debug(fmt.Sprintf("%s,%v", "start redis service...", time.Now()))
	redisConf := redisService.config.RedisServiceConfig.DBMap[redisService.config.RedisServiceConfig.DB]
	if len(redisConf.Password) > 0 {
		redisService.Client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
			Username: redisConf.User,
			Password: redisConf.Password,
			DB:       redisConf.Database,
			PoolSize: 100,
		})
	} else {
		redisService.Client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
			DB:       redisConf.Database,
			PoolSize: 100,
		})
	}
	redisService.State = true
}

func (redisService *RedisService) Stop() {
	fmt.Println("stop redis service...", time.Now())
	redisService.logger.Debug(fmt.Sprintf("%s,%v", "stop redis service...", time.Now()))
	redisService.cancel()
	_ = redisService.Client.Close()
	redisService.Client = nil
	redisService.State = false
}

func (redisService *RedisService) NewRedis(db string) (redisClient *redis.Client) {
	redisConf := redisService.config.RedisServiceConfig.DBMap[db]
	if len(redisConf.Password) > 0 {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
			Username: redisConf.User,
			Password: redisConf.Password,
			DB:       redisConf.Database,
			PoolSize: 100,
		})
	} else {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
			DB:       redisConf.Database,
			PoolSize: 100,
		})
	}
	return
}

func (redisService *RedisService) TryLock(key, val string, timeout int) (lock *Lock, ok bool, err error) {
	lock = &Lock{key: key, val: val, timeout: timeout}
	ok, err = redisService.SetNX(redisService.context, key, val, time.Duration(timeout)*time.Second).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if !ok {
		lock = nil
	}
	return
}

func (redisService *RedisService) AddLockTimeout(lock *Lock, extraTime int) (ok bool, err error) {
	ttl, err := redisService.TTL(redisService.context, lock.key).Result()
	if err != nil {
		return false, err
	}
	if ttl > 0 {
		err = redisService.Set(redisService.context, lock.key, lock.val, ttl+time.Duration(extraTime)*time.Second).Err()
		if err != redis.Nil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func (redisService *RedisService) Unlock(lock *Lock) error {
	return redisService.Del(redisService.context, lock.key).Err()
}
