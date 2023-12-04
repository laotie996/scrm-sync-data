app_name: scrm-task
server_node_num: 1
logger_service_config:
  level: -1
  output_path: ./logs/
  file_name: scrm-task.log
  max_age: 168
  rotate_time_level: 1
  rotate_time: 0
clickhouse_service_config:
  enable: true
  db: statistics
  db_map:
    statistics:
      user: default
      password:
      host: 192.168.0.101
      port: 9000
      database: statistics
      dialect: mysql
      db_source: root:@tcp(192.168.0.101:9000)/statistics?charset=utf8mb4&parseTime=True&loc=Local
      db_debug: true #是否输出数据库日志
redis_service_config:
  enable: true
  db: tg  #项目默认选择db名称
  db_map: #数据库配置
    sys:
      host: 192.168.0.101
      port: 6379
      database: 0
      prefix: test
    tg: #消息服务数据库
      host: 192.168.0.101
      port: 23479
      password:
      database: 0
      prefix: test
mongo_db_service_config:
  enable: true
  data_base: scrm-task  #项目默认选择db名称
  db_map: #数据库配置
    scrm-task: #消息服务数据库
      enable: true
      user: root
      password: 18nY67fvR5
      host: 8.134.166.200
      port: 27017
      database: scrm-task
      db_source: mongodb://root:18nY67fvR5@8.134.166.200:27017/scrm-task?connect=direct
      listen_db_source: mongodb://root:18nY67fvR5@8.134.166.200:27017/scrm-task?connect=direct
      max_idle_count: 10
      max_open_count: 100
    tg-cst-app: #消息服务数据库
      enable: true
      user: cst-admin
      password: 18nY67fvR5
      host: 8.134.166.200
      port: 27017
      database: tg-cst-app
      db_source: mongodb://cst-admin:tg-cstZ8Gs4AFw@8.134.166.200:27017/tg-cst-app?connect=direct
      listen_db_source: mongodb://cst-admin:tg-cstZ8Gs4AFw@8.134.166.200:27017/tg-cst-app?connect=direct
      max_idle_count: 10
      max_open_count: 100
