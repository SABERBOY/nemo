package controller

import (
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/service"
	"netease-kit/nemo/pkg/response"

	"github.com/gin-gonic/gin"
)

type KtvController struct {
	KtvService *service.KtvService
}

func NewKtvController() *KtvController {
	return &KtvController{
		KtvService: service.NewKtvService(),
	}
}

func (ctrl *KtvController) SingStart(c *gin.Context) {
	var param dto.SingParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	if err := ctrl.KtvService.SingStart(userUuid, param); err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (ctrl *KtvController) GetSingInfo(c *gin.Context) {
	var param dto.SingInfoParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	info, err := ctrl.KtvService.GetSingInfo(param.RoomUuid)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, info)
}

func (ctrl *KtvController) SingControl(c *gin.Context) {
	var param dto.SingActionParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	if err := ctrl.KtvService.SingControl(userUuid, param); err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (ctrl *KtvController) ChorusInvite(c *gin.Context) {
	var param dto.ChorusInviteParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	res, err := ctrl.KtvService.ChorusInvite(userUuid, param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, res)
}

func (ctrl *KtvController) ChorusJoin(c *gin.Context) {
	var param dto.JoinChorusParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	res, err := ctrl.KtvService.JoinChorus(userUuid, param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, res)
}

func (ctrl *KtvController) ChorusCancel(c *gin.Context) {
	var param dto.CancelChorusParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	res, err := ctrl.KtvService.CancelChorus(userUuid, param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, res)
}

func (ctrl *KtvController) ChorusReady(c *gin.Context) {
	var param dto.FinishChorusReadyParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	res, err := ctrl.KtvService.ChorusReady(userUuid, param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, res)
}

func (ctrl *KtvController) AbandonSing(c *gin.Context) {
	var param dto.AbandonSingParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}
	userUuid := c.GetHeader("userUuid")
	if err := ctrl.KtvService.AbandonSing(userUuid, param); err != nil {
		response.Failed(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}
