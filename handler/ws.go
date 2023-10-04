package handler

import (
	// "github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"github.com/jian0209/task-monitor-service/utils"
	"github.com/polevpn/anyvalue"
	"github.com/polevpn/elog"
)

type WsHandler struct{}

func (h *WsHandler) OnConnected(c *gin.Context) {
	// create a new connection
	_, err := utils.NewConnectionPool(eventCallBack)
	if err != nil {
		elog.Error("create connection pool fail", err)
		utils.Response500(c, utils.RESP_INTERNAL_ERROR_CODE)
		return
	}
	utils.Response200(c, utils.RESP_SUCCESS_CODE)
}

// func (h *WsHandler) sendNotification(notification string) error {
// 	return nil
// }

func eventCallBack(av *anyvalue.AnyValue) {
	json, _ := av.EncodeJson()
	// send to server
	elog.Info(string(json))
}
