package handler

import (
	"strings"
	"time"

	"github.com/jian0209/task-monitor-service/utils"
)

// set user token into redis
func SetUserToken(username string, userOrigin string) (string, error) {
	redis := utils.RedisClient

	tokenKey := utils.TOKEN_KEY + utils.RandomString(24)
	tokenTimeLeft := time.Duration(utils.TOKEN_TIME_LEFT) * time.Second
	//save token
	return strings.Split(tokenKey, ":")[2], redis.Set(tokenKey, userOrigin, tokenTimeLeft).Err()
}

// get user token from redis
func GetUserToken(token string) (string, error) {
	redis := utils.RedisClient

	tokenKey := utils.TOKEN_KEY + token
	return redis.Get(tokenKey).Result()
}
