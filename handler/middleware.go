package handler

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jian0209/task-monitor-service/dao"
	"github.com/jian0209/task-monitor-service/utils"
	"github.com/polevpn/elog"
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

// check user token from redis
func CheckLogin(c *gin.Context) {
	token := c.GetHeader("token")

	if token == "" {
		elog.Error("token is empty")
		utils.Response400(c, utils.RESP_TOKEN_INVALID)
		c.Abort()
		return
	}

	userJson, err := utils.RedisClient.Get(utils.TOKEN_KEY + token).Result()

	if err != nil {
		elog.Error(err)
		utils.Response400(c, utils.RESP_TOKEN_INVALID)
		c.Abort()
		return
	}

	user := dao.User{}
	err = json.Unmarshal([]byte(userJson), &user)

	if err != nil {
		elog.Error(err)
		utils.Response400(c, utils.RESP_CHECK_LOGIN_FAILED)
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Set("token", token)
	c.Next()
}
