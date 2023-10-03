package utils

import (
	"github.com/go-redis/redis/v7"
	"github.com/polevpn/anyvalue"
	"gorm.io/gorm"
)

var Config *anyvalue.AnyValue
var DBClient *gorm.DB
var RedisClient *redis.Client
