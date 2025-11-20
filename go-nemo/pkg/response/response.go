package response

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	DefaultOK        = 200
	DefaultOKMessage = "success"
	DefaultFailed    = -1
)

func Success(c *gin.Context, data interface{}) {
	render(c, DefaultOK, DefaultOKMessage, data)
}

func Failed(c *gin.Context, code int, msg string) {
	render(c, code, msg, nil)
}

func render(c *gin.Context, code int, msg string, data interface{}) {
	startTime, exists := c.Get("startTime")
	var costTime string
	if exists {
		start := startTime.(time.Time)
		costTime = fmt.Sprintf("%dms", time.Since(start).Milliseconds())
	}

	requestId := c.GetString("requestId")

	response := gin.H{
		"code":      code,
		"requestId": requestId,
		"costTime":  costTime,
	}

	if code == DefaultOK {
		if data != nil {
			response["data"] = data
		}
	} else {
		response["msg"] = msg
	}

	c.JSON(http.StatusOK, response)
}
