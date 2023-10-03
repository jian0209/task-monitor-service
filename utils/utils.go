package utils

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polevpn/anyvalue"
)

func LoadConfig(config string) (*anyvalue.AnyValue, error) {
	dataBytes, err := os.ReadFile(config)

	if err != nil {
		return nil, err
	}
	return anyvalue.NewFromYaml(dataBytes)
}

func Response200(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": RESP_SUCCESS_CODE,
		"msg":  RESPONSE_TEXT[RESP_SUCCESS_CODE],
		"data": data,
	})
}

func Response400(c *gin.Context, code int) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": code,
		"msg":  RESPONSE_TEXT[code],
		"data": nil,
	})
}

func Response500(c *gin.Context, code int) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": code,
		"msg":  RESPONSE_TEXT[code],
		"data": nil,
	})
}

func RandomString(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var random = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	random.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterBytes[random.Intn(len(letterBytes))]
	}
	return string(b)
}

func GetPasswordHash(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func CheckEmailFormat(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func GetTimeNow() time.Time {
	return time.Now().Local()
}
