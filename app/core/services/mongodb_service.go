package services

import (
	"context"
	"fmt"
	"github.com/dtm-labs/logger"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"scrm-sync-data/app/config"
	"sync"
	"time"
)

type MongoDBService struct {
	context        context.Context
	cancel         context.CancelFunc
	config         *config.Config //配置
	logger         *LoggerService
	connectionPool *MongoDBConnectionPool
	State          bool //服务状态
}

type MongoDBConnectionPool struct {
	maxIdleCount    int
	maxOpenCount    int
	idleConnections map[*mongo.Client]time.Time
	idleTimeout     time.Duration
	requestChan     chan chan *mongo.Client
	mutex           sync.Mutex
}

func (mongoDBConnectionPool *MongoDBConnectionPool) get(logger *LoggerService) *mongo.Client {
	connectionChan := make(chan *mongo.Client)
	mongoDBConnectionPool.mutex.Lock()
	defer mongoDBConnectionPool.mutex.Unlock()
	if len(mongoDBConnectionPool.idleConnections) > 0 {
		logger.Debugf("mongodb connection pool size: %d", len(mongoDBConnectionPool.idleConnections))
		for connection, lastActivity := range mongoDBConnectionPool.idleConnections {
			if time.Since(lastActivity) < mongoDBConnectionPool.idleTimeout {
				delete(mongoDBConnectionPool.idleConnections, connection)
				close(connectionChan)
				return connection
			}
		}
	}
	mongoDBConnectionPool.requestChan <- connectionChan
	connection := <-connectionChan
	close(connectionChan)
	return connection
}

func (mongoDBConnectionPool *MongoDBConnectionPool) put(connection *mongo.Client, logger *LoggerService) {
	mongoDBConnectionPool.mutex.Lock()
	defer mongoDBConnectionPool.mutex.Unlock()
	if len(mongoDBConnectionPool.idleConnections) < mongoDBConnectionPool.maxIdleCount {
		mongoDBConnectionPool.idleConnections[connection] = time.Now()
		logger.Debug("put connection to mongodb connection pool successfully")
	} else {
		_ = connection.Disconnect(context.TODO())
	}
}

func (mongoDBConnectionPool *MongoDBConnectionPool) manage(mongoDB *MongoDBService, logger *LoggerService) {
	ticker := time.NewTicker(mongoDBConnectionPool.idleTimeout)
	for {
		select {
		case connChan := <-mongoDBConnectionPool.requestChan:
			clientOptions := options.Client().ApplyURI(mongoDB.config.MongoDBServiceConfig.DBMap[mongoDB.config.MongoDBServiceConfig.DB].DBSource).SetMaxPoolSize(0)
			client, err := mongo.Connect(mongoDB.context, clientOptions)
			if err != nil {
				logger.Errorf("connect failed:%v", err.Error())
				return
			}
			err = client.Ping(context.TODO(), nil)
			if err != nil {
				logger.Errorf("ping failed:%v", err.Error())
				continue
			}
			connChan <- client
		case <-ticker.C:
			for connection, lastActivity := range mongoDBConnectionPool.idleConnections {
				if time.Since(lastActivity) > mongoDBConnectionPool.idleTimeout {
					logger.Debugf("mongodb connection idle timeout:%+v", connection)
					delete(mongoDBConnectionPool.idleConnections, connection)
					_ = connection.Disconnect(context.TODO())
				}
			}
		}
	}
}

func NewMongoDBConnectionPool(maxIdleCount int, maxOpenCount int, mongoDB *MongoDBService, logger *LoggerService) *MongoDBConnectionPool {
	mongoDBConnectionPool := MongoDBConnectionPool{
		maxIdleCount:    maxIdleCount,
		maxOpenCount:    maxOpenCount,
		idleConnections: make(map[*mongo.Client]time.Time),
		idleTimeout:     time.Duration(10) * time.Second,
		requestChan:     make(chan chan *mongo.Client),
		mutex:           sync.Mutex{},
	}
	go mongoDBConnectionPool.manage(mongoDB, logger)
	return &mongoDBConnectionPool
}

func (mongoDBService *MongoDBService) Init(parentContext context.Context, config *config.Config, logger *LoggerService) {
	mongoDBService.config = config
	mongoDBService.logger = logger
	mongoDBService.State = false
	mongoDBService.context, mongoDBService.cancel = context.WithCancel(parentContext)
	mongoDBService.Start()
}

func (mongoDBService *MongoDBService) Start() {
	fmt.Println("start mongodb service...", time.Now())
	mongoDBService.logger.Debug(fmt.Sprintf("%s,%v", "start mongodb service...", time.Now()))
	mongoDBConnectionPool := NewMongoDBConnectionPool(
		mongoDBService.config.MongoDBServiceConfig.DBMap[mongoDBService.config.MongoDBServiceConfig.DB].MaxIdleCount,
		mongoDBService.config.MongoDBServiceConfig.DBMap[mongoDBService.config.MongoDBServiceConfig.DB].MaxOpenCount,
		mongoDBService,
		mongoDBService.logger,
	)
	mongoDBService.connectionPool = mongoDBConnectionPool
	mongoDBService.State = true
}

func (mongoDBService *MongoDBService) Stop() {
	fmt.Println("stop mongodb service...", time.Now())
	mongoDBService.logger.Debug(fmt.Sprintf("%s,%v", "stop mongodb service...", time.Now()))
	mongoDBService.cancel()
	mongoDBService.State = false
}

func (mongoDBService *MongoDBService) NewConnection() *mongo.Client {
	connection := mongoDBService.connectionPool.get(mongoDBService.logger)
	return connection
}

func (mongoDBService *MongoDBService) ReleaseConnection(connection *mongo.Client) {
	mongoDBService.connectionPool.put(connection, mongoDBService.logger)
	return
}

func (mongoDBService *MongoDBService) NewDB(dbName string) (*mongo.Client, error) {
	fmt.Printf("start mongodb:%s service...%s\n", dbName, time.Now().Format(time.DateTime))
	mongoDBService.logger.Debug(fmt.Sprintf("%s,%v", "start mongodb service...", time.Now()))
	dbConfig, ok := mongoDBService.config.MongoDBServiceConfig.DBMap[dbName]
	if !ok || dbConfig.DBSource == "" {
		return nil, errors.New("配置不存在")
	}
	clientOptions := options.Client().ApplyURI(dbConfig.DBSource).SetMaxPoolSize(0)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.Errorf("connect failed:%v", err.Error())
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Errorf("ping failed:%v", err.Error())
		return nil, err
	}
	return client, nil
}
