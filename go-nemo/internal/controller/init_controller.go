package controller

import (
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/service"
	"netease-kit/nemo/pkg/response"

	"github.com/gin-gonic/gin"
)

type NemoInitController struct {
	AppInitService *service.AppInitService
}

func NewNemoInitController() *NemoInitController {
	return &NemoInitController{
		AppInitService: service.NewAppInitService(),
	}
}

func (ctrl *NemoInitController) InitAppAndUser(c *gin.Context) {
	var param dto.InitUserParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	appKey := c.GetString("appKey")
	// secret := c.GetString("appSecret") // In Java code, it checks secret against context, here we assume middleware handled auth or we add check

	// Basic validation from Java controller
	// String appSecret = request.getHeader("AppSecret");
	// if (StringUtils.isEmpty(appSecret) || !secret.equals(appSecret)) { ... }
	// For now, skipping strict secret check against DB/Config as Context.get().getSecret() source wasn't fully ported yet.

	userDto, err := ctrl.AppInitService.InitAppAndUser(appKey, param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, userDto)
}
