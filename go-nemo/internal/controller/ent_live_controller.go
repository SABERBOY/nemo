package controller

import (
	"netease-kit/nemo/internal/service"
	"netease-kit/nemo/pkg/response"

	"github.com/gin-gonic/gin"
)

type EntLiveController struct {
	EntLiveService *service.EntLiveService
}

func NewEntLiveController() *EntLiveController {
	return &EntLiveController{
		EntLiveService: service.NewEntLiveService(),
	}
}

type CreateLiveParam struct {
	LiveTopic string `json:"liveTopic"`
	Cover     string `json:"cover"`
	LiveType  int    `json:"liveType"`
}

func (ctrl *EntLiveController) CreateLive(c *gin.Context) {
	var param CreateLiveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	// Assuming userUuid is extracted from token/header in a real scenario
	// For demo, we might need to pass it or get from context if auth middleware set it
	// userUuid := c.GetString("userUuid")
	// Using a dummy for now as auth middleware isn't fully fleshed out with token parsing
	userUuid := c.GetHeader("userUuid")
	if userUuid == "" {
		response.Failed(c, 401, "User not identified")
		return
	}

	record, err := ctrl.EntLiveService.CreateLive(userUuid, param.LiveTopic, param.Cover, param.LiveType)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, record)
}

type CloseLiveParam struct {
	RoomArchiveId string `json:"roomArchiveId"`
}

func (ctrl *EntLiveController) CloseLive(c *gin.Context) {
	var param CloseLiveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	userUuid := c.GetHeader("userUuid")
	if userUuid == "" {
		response.Failed(c, 401, "User not identified")
		return
	}

	err := ctrl.EntLiveService.CloseLive(userUuid, param.RoomArchiveId)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *EntLiveController) GetLiveList(c *gin.Context) {
	type ListParam struct {
		LiveType int `form:"liveType"`
	}
	var param ListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	userUuid := c.GetHeader("userUuid")
	// userUuid can be empty for list

	list, err := ctrl.EntLiveService.GetLiveList(param.LiveType, userUuid)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, list)
}

func (ctrl *EntLiveController) GetLiveInfo(c *gin.Context) {
	type InfoParam struct {
		LiveRecordId uint64 `form:"liveRecordId"`
	}
	var param InfoParam
	if err := c.ShouldBindQuery(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	info, err := ctrl.EntLiveService.GetLiveInfo(param.LiveRecordId)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	if info == nil {
		response.Failed(c, 404, "Live record not found")
		return
	}

	response.Success(c, info)
}
