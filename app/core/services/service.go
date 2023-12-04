package services

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"runtime"
	"scrm-sync-data/app/config"
)

// Service
// @Description: 核心服务组件
type Service interface {
	Start()
	Stop()
}

type Services struct {
	Context        context.Context
	Cancel         context.CancelFunc
	Config         *config.Config
	Logger         *LoggerService
	Network        *NetworkService
	Format         *FormatService
	Snowflake      *snowflake.Node
	Crypto         *CryptoService
	DelayTaskQueue *DelayTaskQueueService
	Clickhouse     *ClickhouseService
	Redis          *RedisService
	MongoDB        *MongoDBService
}

// Init
// @Description:核心模块服务初始化
// @receiver services
func (services *Services) Init(config *config.Config) {
	services.Context, services.Cancel = context.WithCancel(context.TODO())
	services.Config = config
	services.Logger = new(LoggerService)
	services.Logger.Init(services.Context, services.Config)
	services.Network = new(NetworkService)
	services.Network.Init(services.Context, services.Config)
	services.Format = new(FormatService)
	services.Format.Init(services.Context, services.Config, services.Logger)
	services.Snowflake, _ = snowflake.NewNode(services.Config.ServerNodeNum)
	services.Crypto = new(CryptoService)
	services.Crypto.Init(services.Context, services.Config, services.Logger)
	services.DelayTaskQueue = new(DelayTaskQueueService)
	services.DelayTaskQueue.Init(services.Context, services.Config, services.Logger)
	services.Clickhouse = new(ClickhouseService)
	services.Clickhouse.Init(services.Context, services.Config, services.Logger)
	if config.RedisServiceConfig.Enable {
		services.Redis = new(RedisService)
		services.Redis.Init(services.Context, services.Config, services.Logger)
	}
	if config.MongoDBServiceConfig.Enable {
		services.MongoDB = new(MongoDBService)
		services.MongoDB.Init(services.Context, services.Config, services.Logger)
	}
}

// GetServices
// @Description: 获取核心模块所有服务实例
// @receiver services
// @return *Services
func (services *Services) GetServices() *Services {
	return services
}

// GetService
// @Description: 根据服务名称获取核心模块单个服务实例
// @receiver services
// @return Service
func (services *Services) GetService(serviceName string) Service {
	for {
		select {
		case <-services.Context.Done():
			return nil
		default:
			switch serviceName {
			case "logger":
				return services.Logger
			case "network":
				return services.Network
			case "format":
				return services.Format
			case "crypto":
				return services.Crypto
			case "delay-task-queue":
				return services.DelayTaskQueue
			case "clickhouse":
				return services.Clickhouse
			case "redis":
				return services.Redis
			case "mongodb":
				return services.MongoDB
			}
		}
	}
}

// GetLogger
// @Description: 获取核心模块日志服务实例
// @receiver services
// @return Service
func (services *Services) GetLogger() Service {
	for {
		select {
		case <-services.Context.Done():
			return nil
		default:
			return services.Logger
		}
	}
}

// GetNetwork
// @Description: 获取核心模块网络服务实例
// @receiver services
// @return Service
func (services *Services) GetNetwork() Service {
	for {
		select {
		case <-services.Context.Done():
			return nil
		default:
			return services.Network
		}
	}
}

// GetCrypto
// @Description: 获取核心模块加密和解密服务实例
// @receiver services
// @return Service
func (services *Services) GetCrypto() Service {
	for {
		select {
		case <-services.Context.Done():
			return nil
		default:
			return services.Crypto
		}
	}
}

// GetDelayTaskQueue
// @Description: 获取核心模块延时任务队列服务实例
// @receiver services
// @return Service
func (services *Services) GetDelayTaskQueue() Service {
	for {
		select {
		case <-services.Context.Done():
			return nil
		default:
			return services.DelayTaskQueue
		}
	}
}

// GetClickhouse
// @Description: 获取核心模块clickhouse数据库服务实例
// @receiver services
// @return Service
func (services *Services) GetClickhouse() Service {
	for {
		select {
		case <-services.Context.Done():
			return nil
		default:
			return services.Clickhouse
		}
	}
}

// GetRedis
// @Description: 获取核心模块Redis服务实例
// @receiver services
// @return Service
func (services *Services) GetRedis() Service {
	for {
		select {
		case <-services.Context.Done():
			return nil
		default:
			return services.Redis
		}
	}
}

// GetMongoDB
// @Description: 获取核心模块MongoDB服务实例
// @receiver services
// @return Service
func (services *Services) GetMongoDB() Service {
	for {
		select {
		case <-services.Context.Done():
			return nil
		default:
			return services.MongoDB
		}
	}
}

func (services *Services) SafeGo(logger *LoggerService, f func(args ...interface{})) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<16)
			n := runtime.Stack(buf, true)
			logger.Errorf("%v %s", r, string(buf[:n]))
			panic(r)
		}
	}()
	f()
}
