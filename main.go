package main

import (
	"flag"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jian0209/task-monitor-service/handler"
	"github.com/jian0209/task-monitor-service/utils"
	"github.com/jian0209/task-monitor-service/websocket"
	"github.com/polevpn/elog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var configFile string
var err error

func init() {
	flag.StringVar(&configFile, "configFile", "config.yaml", "config file path")
}

func initConfig() error {
	utils.Config, err = utils.LoadConfig(configFile)
	return err
}

func initMysql() error {
	dsn := utils.Config.Get("mysql.dsn").AsStr("root:asdQWE123#@tcp(127.0.0.1:3306)/task_monitor?charset=utf8mb4&parseTime=True&loc=Local")
	namingStrategy := schema.NamingStrategy{SingularTable: true}
	utils.DBClient, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{NamingStrategy: namingStrategy})
	sqlDB, err := utils.DBClient.DB()
	sqlDB.SetMaxIdleConns(utils.Config.Get("mysql.max_idle_conn").AsInt(10))
	sqlDB.SetMaxOpenConns(utils.Config.Get("mysql.max_open_conn").AsInt(100))
	sqlDB.SetConnMaxLifetime(time.Hour)
	return err
}

func initRedis() error {
	utils.RedisClient = redis.NewClient(&redis.Options{
		Addr:         utils.Config.Get("redis.addr").AsStr("127.0.0.1:6379"),
		Password:     utils.Config.Get("redis.password").AsStr(""),
		DB:           utils.Config.Get("redis.db").AsInt(0),
		PoolSize:     utils.Config.Get("redis.pool_size").AsInt(500),
		MinIdleConns: utils.Config.Get("redis.min_idle_conn").AsInt(10),
	})
	return nil
}

func initGinRouter() (*gin.Engine, error) {
	userHandler := &handler.UserHandler{}

	gin.SetMode(utils.Config.Get("gin.mode").AsStr("debug"))
	ginLogFile := utils.Config.Get("gin.log_mode").AsStr("console")

	if ginLogFile == "console" {
		gin.DefaultWriter = io.MultiWriter(os.Stderr)
	} else {
		f, _ := os.Create(ginLogFile)
		gin.DefaultWriter = io.MultiWriter(f)
	}

	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/ws", func(c *gin.Context) {
		handler := websocket.NewRequestHandler()
		server := websocket.NewHttpServer(handler)
		err := server.Listen("127.0.0.1:9011")

		if err != nil {
			elog.Error("start server fail:", err)
			return
		}
	})

	r.POST("/api/user/register", userHandler.Register)
	r.POST("/api/user/login", userHandler.Login)

	return r, nil
}

func main() {
	flag.Parse()
	elog.SetLogToStderr(true)
	defer elog.Flush()

	elog.Info("init config")
	err = initConfig()
	if err != nil {
		elog.Error("init config failed", err)
		return
	}

	elog.Info("init mysql")
	err = initMysql()
	if err != nil {
		elog.Error("init mysql fail:", err)
		return
	}

	elog.Info("init redis")
	err = initRedis()
	if err != nil {
		elog.Error("init to redis fail:", err)
		return
	}

	elog.Info("init gin router")
	r, _ := initGinRouter()
	elog.Error(r.Run(utils.Config.Get("listen").AsStr(":8080")))
}
