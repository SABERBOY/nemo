package controller

import (
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/service"
	"netease-kit/nemo/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MusicPlayerController struct {
	MusicPlayService *service.MusicPlayService
}

func NewMusicPlayerController() *MusicPlayerController {
	return &MusicPlayerController{
		MusicPlayService: service.NewMusicPlayService(),
	}
}

func (ctrl *MusicPlayerController) GetPlayMusicInfo(c *gin.Context) {
	liveRecordIdStr := c.Query("liveRecordId")
	liveRecordId, err := strconv.ParseUint(liveRecordIdStr, 10, 64)
	if err != nil {
		response.Failed(c, 400, "Invalid liveRecordId")
		return
	}

	info, err := ctrl.MusicPlayService.GetPlayMusicInfo(liveRecordId)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, info)
}

func (ctrl *MusicPlayerController) MusicAction(c *gin.Context) {
	var param dto.MusicActionParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	if err := ctrl.MusicPlayService.MusicAction(userUuid, param); err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (ctrl *MusicPlayerController) MusicReady(c *gin.Context) {
	var param dto.MusicReadyParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	if err := ctrl.MusicPlayService.MusicReady(param.LiveRecordId, param.OrderId, userUuid); err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}
