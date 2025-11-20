package controller

import (
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/service"
	"netease-kit/nemo/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SongController struct {
	KtvService *service.KtvService
}

func NewSongController() *SongController {
	return &SongController{
		KtvService: service.NewKtvService(),
	}
}

func (ctrl *SongController) OrderSong(c *gin.Context) {
	var param dto.OrderSongParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	res, err := ctrl.KtvService.OrderSong(userUuid, param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, res)
}

func (ctrl *SongController) SwitchSong(c *gin.Context) {
	var param dto.SwitchSongParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	if err := ctrl.KtvService.SwitchSong(param.LiveRecordId, userUuid, param.CurrentOrderId, param.Attachment); err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (ctrl *SongController) GetOrderSongs(c *gin.Context) {
	liveRecordIdStr := c.Query("liveRecordId")
	liveRecordId, err := strconv.ParseUint(liveRecordIdStr, 10, 64)
	if err != nil {
		response.Failed(c, 400, "Invalid liveRecordId")
		return
	}

	list, err := ctrl.KtvService.GetOrderSongs(liveRecordId)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, list)
}

func (ctrl *SongController) CancelOrderSong(c *gin.Context) {
	var param dto.OrderParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	if err := ctrl.KtvService.CancelOrderSong(userUuid, param.LiveRecordId, param.OrderId); err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (ctrl *SongController) SongSetTop(c *gin.Context) {
	var param dto.OrderParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	if err := ctrl.KtvService.SongSetTop(userUuid, param.LiveRecordId, param.OrderId); err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (ctrl *SongController) CleanUserOrderSongs(c *gin.Context) {
	var param dto.CleanOrderSongParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	if err := ctrl.KtvService.CleanUserOrderSongs(param.LiveRecordId, userUuid); err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (ctrl *SongController) GetMusicToken(c *gin.Context) {
	// Mock token
	response.Success(c, map[string]interface{}{
		"accessToken": "mock_token",
		"expiresIn":   21600,
	})
}
