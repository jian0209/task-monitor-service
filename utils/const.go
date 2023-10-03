package utils

const (
	RESP_SUCCESS_CODE        = 0
	RESP_PARAM_ERROR_CODE    = 1000
	RESP_USER_NOT_FOUND_CODE = 1001
	RESP_USER_EMAIL_ERROR    = 1002
	RESP_USER_PASSWORD_ERROR = 1003
	RESP_USER_EXIST_CODE     = 1004
	RESP_INTERNAL_ERROR_CODE = 1099
)

var RESPONSE_TEXT = map[int]string{
	0:    "success",
	1000: "param error",
	1001: "invalid username or password",
	1002: "invalid email",
	1003: "invalid password",
	1004: "user exist",
	1099: "internal error",
}

const (
	SERVICE_NAME   = "task-monitor-service"
	WEBSOCKET_NAME = "ws://127.0.0.1:9010/ws"
)

const (
	PASSWORD_SALT       = "sLx2j3SfvE&"
	PASSWORD_MIN_LENGTH = 6
)

const (
	TOKEN_TIME_LEFT = 3600 * 24 * 7
	TOKEN_KEY       = SERVICE_NAME + ":token:"
)

const (
	TCP_WRITE_BUFFER_SIZE   = 524288
	TCP_READ_BUFFER_SIZE    = 524288
	CH_WEBSOCKET_WRITE_SIZE = 2000
	TRAFFIC_LIMIT_INTERVAL  = 10
)
