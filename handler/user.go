package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jian0209/task-monitor-service/dao"
	"github.com/jian0209/task-monitor-service/utils"
	"github.com/polevpn/elog"
)

type UserHandler struct{}

/*
@ Request method: POST
@ Request body:
@ username: string
@ password: string
@ email: string
@ response:
@ code: int
@ msg: string
@ data: nil
*/
func (uh *UserHandler) Register(c *gin.Context) {
	// bind request
	var req utils.RegisterReq
	if err := c.BindJSON(&req); err != nil {
		elog.Error(err.Error())
		utils.Response400(c, utils.RESP_PARAM_ERROR_CODE)
		return
	}

	// check email format
	if !utils.CheckEmailFormat(req.Email) {
		elog.Error("invalid email: " + req.Email)
		utils.Response400(c, utils.RESP_USER_EMAIL_ERROR)
		return
	}

	// check password length
	if len(req.Password) < utils.PASSWORD_MIN_LENGTH {
		elog.Error("invalid password length: " + string(len(req.Password)))
		utils.Response400(c, utils.RESP_USER_PASSWORD_ERROR)
		return
	}

	// set password with salt
	randomString := utils.RandomString(10)
	req.Password = utils.GetPasswordHash(utils.PASSWORD_SALT + req.Password + randomString)

	db := utils.DBClient
	userDao := dao.NewUserDAO(db)
	// check user exist
	user, err := userDao.GetByUsername(req.Username)
	if err != nil {
		elog.Error(err.Error())
		utils.Response500(c, utils.RESP_INTERNAL_ERROR_CODE)
		return
	}

	if user != nil {
		elog.Warn("user exist: " + req.Username)
		utils.Response400(c, utils.RESP_USER_EXIST_CODE)
		return
	}

	// set user info
	userInfo := make(map[string]interface{})
	userInfo["username"] = req.Username
	userInfo["password"] = req.Password
	userInfo["email"] = req.Email
	userInfo["app_secret"] = randomString
	userInfo["status"] = 1

	// insert new user into db
	if err := userDao.Add(userInfo); err != nil {
		elog.Error(err.Error())
		utils.Response500(c, utils.RESP_INTERNAL_ERROR_CODE)
		return
	}

	utils.Response200(c, nil)
}

/*
@ Request method: POST
@ Request body:
@ username: string
@ password: string
@ response:
@ code: int
@ msg: string
@ data: map[string]interface{}
@ data.token: string
@ data.username: string
@ data.email: string
@ data.status: int
@ data.updated_at: string
@ data.created_at: string
*/
func (uh *UserHandler) Login(c *gin.Context) {
	var req utils.LoginReq
	if err := c.BindJSON(&req); err != nil {
		elog.Error(err.Error())
		utils.Response400(c, utils.RESP_PARAM_ERROR_CODE)
		return
	}

	// get user info from db by username first, and get the app secret from the username
	// then set the password with salt and app secret
	db := utils.DBClient
	userDao := dao.NewUserDAO(db)
	result, err := userDao.GetByUsername("jian0209")
	if err != nil {
		elog.Error(err.Error())
		utils.Response500(c, utils.RESP_INTERNAL_ERROR_CODE)
		return
	}

	if result == nil {
		elog.Warn("user not found: " + req.Username)
		utils.Response400(c, utils.RESP_USER_NOT_FOUND_CODE)
		return
	}

	userAppSecret := result.App_secret
	req.Password = utils.GetPasswordHash(utils.PASSWORD_SALT + req.Password + userAppSecret)

	// compare password
	if req.Password != result.Password {
		elog.Warn("invalid password: " + req.Username)
		utils.Response400(c, utils.RESP_USER_NOT_FOUND_CODE)
		return
	}

	userOrigin, _ := json.Marshal(dao.ConvertToUserRedis(result))

	// set token
	token, err := SetUserToken(result.Username, string(userOrigin))
	if err != nil {
		elog.Error(err.Error())
		utils.Response500(c, utils.RESP_INTERNAL_ERROR_CODE)
		return
	}

	updateInfo := make(map[string]interface{})
	updateInfo["last_login_ip"] = c.ClientIP()
	updateInfo["last_login_time"] = utils.GetTimeNow()

	// update user info
	if err := userDao.Update(result.Id, updateInfo); err != nil {
		elog.Error(err.Error())
		utils.Response500(c, utils.RESP_INTERNAL_ERROR_CODE)
		return
	}

	returnResult := make(map[string]interface{})

	returnResult["token"] = token
	returnResult["username"] = result.Username
	returnResult["email"] = result.Email
	returnResult["status"] = result.Status
	returnResult["updated_at"] = result.Updated_At
	returnResult["created_at"] = result.Created_At

	utils.Response200(c, returnResult)
}
