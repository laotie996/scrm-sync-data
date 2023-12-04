package services

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"scrm-sync-data/app/config"
	"time"
)

type ClickhouseService struct {
	context context.Context
	cancel  context.CancelFunc
	config  *config.Config
	logger  *LoggerService
	driver.Conn
	error error
	State bool
}

func (clickhouseService *ClickhouseService) Init(parentContext context.Context, config *config.Config, logger *LoggerService) {
	clickhouseService.config = config
	clickhouseService.logger = logger
	clickhouseService.State = false
	clickhouseService.context, clickhouseService.cancel = context.WithCancel(parentContext)
	clickhouseService.Start()
}

func (clickhouseService *ClickhouseService) Start() {
	fmt.Println("start clickhouse service...", time.Now())
	dbConf := clickhouseService.config.ClickhouseServiceConfig.DBMap[clickhouseService.config.ClickhouseServiceConfig.DB]
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("start clickhouse service error,%s\n", r)
			clickhouseService.logger.Errorf("start clickhouse service error,%v", r)
			clickhouseService.Stop()
			return
		}
	}()
	if dbConf.DBDebug {
		var loggerFileName = fmt.Sprintf("%s-clickhouse-db-debug.log", dbConf.Database)
		clickhouseLogger := clickhouseService.logger.NewLogger(loggerFileName)
		clickhouseService.Conn, clickhouseService.error = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%d", dbConf.Host, dbConf.Port)},
			Auth: clickhouse.Auth{
				Database: dbConf.Database,
				Username: dbConf.User,
				Password: dbConf.Password,
			},
			Debug: true,
			Debugf: func(format string, v ...any) {
				clickhouseLogger.Debugf(format, v)
			},
			DialTimeout:      time.Second * 30,
			MaxOpenConns:     dbConf.MaxOpenCount,
			MaxIdleConns:     dbConf.MaxIdleCount,
			ConnMaxLifetime:  time.Duration(10) * time.Minute,
			ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		})
	} else {
		clickhouseService.Conn, clickhouseService.error = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%d", dbConf.Host, dbConf.Port)},
			Auth: clickhouse.Auth{
				Database: dbConf.Database,
				Username: dbConf.User,
				Password: dbConf.Password,
			},
			DialTimeout:      time.Second * 30,
			MaxOpenConns:     dbConf.MaxOpenCount,
			MaxIdleConns:     dbConf.MaxIdleCount,
			ConnMaxLifetime:  time.Duration(10) * time.Minute,
			ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		})
	}
	if clickhouseService.error != nil {
		panic(clickhouseService.error)
	}
	clickhouseService.error = clickhouseService.Conn.Ping(context.TODO())
	if clickhouseService.error != nil {
		panic(clickhouseService.error)
	}
	clickhouseService.logger.Debug(fmt.Sprintf("%s,%v", "start clickhouse service...", time.Now()))
	clickhouseService.State = true
}

func (clickhouseService *ClickhouseService) Stop() {
	fmt.Println("stop clickhouse service...", time.Now())
	clickhouseService.logger.Debug(fmt.Sprintf("%s,%v", "stop clickhouse service...", time.Now()))
	clickhouseService.cancel()
	clickhouseService.State = false
}

func (clickhouseService *ClickhouseService) NewDB(db string) (conn driver.Conn, err error) {
	dbConf := clickhouseService.config.ClickhouseServiceConfig.DBMap[db]
	if dbConf.DBDebug {
		var loggerFileName = fmt.Sprintf("%s-clickhouse-db-debug.log", dbConf.Database)
		clickhouseLogger := clickhouseService.logger.NewLogger(loggerFileName)
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%d", dbConf.Host, dbConf.Port)},
			Auth: clickhouse.Auth{
				Database: dbConf.Database,
				Username: dbConf.User,
				Password: dbConf.Password,
			},
			Debug: true,
			Debugf: func(format string, v ...any) {
				clickhouseLogger.Debugf(format, v)
			},
			DialTimeout:      time.Second * 30,
			MaxOpenConns:     dbConf.MaxOpenCount,
			MaxIdleConns:     dbConf.MaxIdleCount,
			ConnMaxLifetime:  time.Duration(10) * time.Minute,
			ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		})
	} else {
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%d", dbConf.Host, dbConf.Port)},
			Auth: clickhouse.Auth{
				Database: dbConf.Database,
				Username: dbConf.User,
				Password: dbConf.Password,
			},
			DialTimeout:      time.Second * 30,
			MaxOpenConns:     dbConf.MaxOpenCount,
			MaxIdleConns:     dbConf.MaxIdleCount,
			ConnMaxLifetime:  time.Duration(10) * time.Minute,
			ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		})
	}
	if err != nil {
		clickhouseService.logger.Error(err.Error())
		return nil, err
	}
	err = conn.Ping(context.TODO())
	if err != nil {
		clickhouseService.logger.Error(err.Error())
		return nil, err
	}
	return
}
