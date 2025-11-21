package controller

import (
	"io"
	"log"
	"net/http"
	"netease-kit/nemo/internal/config"
	"netease-kit/nemo/internal/service"
	"netease-kit/nemo/pkg/utils"

	"github.com/gin-gonic/gin"
)

type NotifyController struct {
	NotifyService  *service.NotifyService
	ImEventService *service.ImEventService
}

func NewNotifyController(notifyService *service.NotifyService, imEventService *service.ImEventService) *NotifyController {
	return &NotifyController{
		NotifyService:  notifyService,
		ImEventService: imEventService,
	}
}

func (c *NotifyController) Notify(ctx *gin.Context) {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 414})
		return
	}
	body := string(bodyBytes)
	if body == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 414})
		return
	}

	curTime := ctx.GetHeader("CurTime")
	checksum := ctx.GetHeader("CheckSum")
	appSecret := config.AppConfig.Yunxin.Origin.AppSecret

	verifyMD5 := utils.GetMD5(body)
	verifyChecksum := utils.GetCheckSum(verifyMD5, curTime, appSecret)

	if verifyChecksum == checksum {
		c.NotifyService.HandleNotify(body)
		ctx.JSON(http.StatusOK, gin.H{"code": 200})
	} else {
		log.Printf("Bad checksum. Expected: %s, Got: %s", verifyChecksum, checksum)
		ctx.JSON(http.StatusOK, gin.H{"code": 414})
	}
}

func (c *NotifyController) ImEventNotify(ctx *gin.Context) {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400})
		return
	}
	body := string(bodyBytes)
	if body == "" {
		ctx.JSON(http.StatusOK, gin.H{"code": 400})
		return
	}

	appKey := ctx.GetHeader("AppKey")
	if appKey != config.AppConfig.Yunxin.Origin.AppKey {
		log.Printf("AppKey mismatch. Expected: %s, Got: %s", config.AppConfig.Yunxin.Origin.AppKey, appKey)
		ctx.JSON(http.StatusOK, gin.H{"code": 200}) // Ignore
		return
	}

	curTime := ctx.GetHeader("CurTime")
	checksum := ctx.GetHeader("CheckSum")
	appSecret := config.AppConfig.Yunxin.Origin.AppSecret

	verifyMD5 := utils.GetMD5(body)
	verifyChecksum := utils.GetCheckSum(verifyMD5, curTime, appSecret)

	if verifyChecksum != checksum {
		log.Printf("Bad checksum. Expected: %s, Got: %s", verifyChecksum, checksum)
		ctx.JSON(http.StatusOK, gin.H{"code": 400})
		return
	}

	c.ImEventService.HandlerChatroomInOut(body)
	ctx.JSON(http.StatusOK, gin.H{"code": 200})
}
