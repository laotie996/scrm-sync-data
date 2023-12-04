package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ClickhouseDB struct {
	User         string `yaml:"user" desc:"用户名"`
	Password     string `yaml:"password" desc:"用户密码"`
	Host         string `yaml:"host" desc:"主机名"`
	Port         int    `yaml:"port" desc:"主机端口"`
	Database     string `yaml:"database" desc:"数据库名"`
	Dialect      string `yaml:"dialect" desc:"数据库类型"`
	DBSource     string `yaml:"db_source" desc:"数据库源"`
	MaxIdleCount int    `yaml:"max_idle_count" desc:"mongodb连接池最大空闲数"`
	MaxOpenCount int    `yaml:"max_open_count" desc:"mongodb连接池最大打开数"`
	DBDebug      bool   `yaml:"db_debug" desc:"是否输出gorm数据库调试语句"`
}

type RedisDB struct {
	Host     string `yaml:"host" desc:"主机名"`
	Port     int    `yaml:"port" desc:"主机端口"`
	User     string `yaml:"user" desc:"redis用户名"`
	Password string `yaml:"password" desc:"redis鉴权密码"`
	Database int    `yaml:"database" desc:"数据库"`
	Prefix   string `yaml:"prefix" desc:"数据库命名前缀"`
}

type LoggerServiceConfig struct {
	Level           int    `yaml:"level" desc:"日志等级 debug:-1 info:0 warn: 1 error:2 dpanic:3 panic:4 fatal:5" json:"level,omitempty"`
	OutputPath      string `yaml:"output_path" desc:"输出路径" json:"output_path,omitempty"`
	FileName        string `yaml:"file_name" desc:"日志名称" json:"file_name,omitempty"`
	MaxAge          int    `yaml:"max_age" desc:"日志最大保留时间" json:"max_age,omitempty"`
	RotateTimeLevel int    `yaml:"rotate_time_level" desc:"日志分片时间等级 0 自定义时间分片 1 日分片 2 1小时分片 3 1分钟分片" json:"rotate_time_level,omitempty"`
	RotateTime      int    `yaml:"rotate_time" desc:"自定义时间分片时长 单位为:min" json:"rotate_time,omitempty"`
}

type ClickhouseServiceConfig struct {
	DB     string                  `yaml:"db" desc:"clickhouse数据库名"`
	DBMap  map[string]ClickhouseDB `yaml:"db_map" desc:"clickhouse数据库配置集合"`
	Enable bool                    `yaml:"enable" desc:"服务是否开启"`
}

type RedisServiceConfig struct {
	DB     string             `yaml:"db" desc:"redis数据库名"`
	DBMap  map[string]RedisDB `yaml:"db_map" desc:"redis数据库配置集合"`
	Enable bool               `yaml:"enable" desc:"服务是否开启"`
}

type MongoDBServiceConfig struct {
	DB     string             `yaml:"db" desc:"mongo数据库名"`
	DBMap  map[string]MongoDB `yaml:"db_map" desc:"mongo数据库配置集合"`
	Enable bool               `yaml:"enable" desc:"服务是否开启"`
}

type MongoDB struct {
	User           string `yaml:"user" desc:"mongodb鉴权用户名"`
	Password       string `yaml:"password" desc:"mongodb鉴权密码"`
	Host           string `yaml:"host" desc:"主机名"`
	Port           int    `yaml:"port" desc:"主机端口"`
	Database       string `yaml:"database" desc:"mongodb数据库名"`
	DBSource       string `yaml:"db_source" desc:"数据库源"`
	ListenDBSource string `yaml:"listen_db_source" desc:"数据同步监听的数据库源"`
	MaxIdleCount   int    `yaml:"max_idle_count" desc:"mongodb连接池最大空闲数"`
	MaxOpenCount   int    `yaml:"max_open_count" desc:"mongodb连接池最大打开数"`
}

type Config struct {
	AppName                 string                  `yaml:"app_name" desc:"服务器应用名称" json:"appName,omitempty"`
	ServerNodeNum           int64                   `yaml:"server_node_num" desc:"服务器端节点编码" json:"serverNodeNum,omitempty"`
	Development             bool                    `yaml:"development" desc:"项目模式 生产模式:false  开发模式:true" json:"development,omitempty"`
	LoggerServiceConfig     LoggerServiceConfig     `yaml:"logger_service_config" desc:"日志配置" json:"loggerServiceConfig"`
	ClickhouseServiceConfig ClickhouseServiceConfig `yaml:"clickhouse_service_config" desc:"clickhouse数据库服务配置" json:"clickhouseServiceConfig"`
	RedisServiceConfig      RedisServiceConfig      `yaml:"redis_service_config" desc:"redis数据库服务配置" json:"redisServiceConfig"`
	MongoDBServiceConfig    MongoDBServiceConfig    `yaml:"mongo_db_service_config" desc:"mongodb数据库服务配置" json:"mongoDBServiceConfig"`
}

func (config *Config) Load() error {
	configFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return fmt.Errorf("read config file failed:%s\n", err.Error())
	}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		return fmt.Errorf("parse config file failed:%s\n", err.Error())
	}
	return nil
}

func (config *Config) Save() error {
	configBuffer, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("parse config to config buffer failed:%s\n", err.Error())
	}
	err = os.WriteFile("./config/config.yaml", configBuffer, 0666)
	if err != nil {
		return fmt.Errorf("write config to file failed:%s\n", err.Error())
	}
	return nil
}
