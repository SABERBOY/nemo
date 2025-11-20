package controller

import (
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/service"
	"netease-kit/nemo/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SocialChatController struct {
	SocialChatService *service.SocialChatService
}

func NewSocialChatController() *SocialChatController {
	return &SocialChatController{
		SocialChatService: service.NewSocialChatService(),
	}
}

type ReporterParam struct {
	UserUuid string `json:"userUuid"`
	DeviceId string `json:"deviceId"`
}

func (ctrl *SocialChatController) Reporter(c *gin.Context) {
	var param ReporterParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	appKey := c.GetString("appKey")
	err := ctrl.SocialChatService.Reporter(appKey, param.UserUuid, param.DeviceId)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *SocialChatController) GetOnLineUser(c *gin.Context) {
	appKey := c.GetString("appKey")
	pageNumStr := c.Query("pageNum")
	pageSizeStr := c.Query("pageSize")

	pageNum, _ := strconv.Atoi(pageNumStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	users, err := ctrl.SocialChatService.GetOnLineUser(appKey, pageNum, pageSize)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, users)
}

func (ctrl *SocialChatController) UserReward(c *gin.Context) {
	var param dto.UserRewardDto
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	err := ctrl.SocialChatService.UserReward(param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}
