package controller

import (
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/service"
	"netease-kit/nemo/pkg/response"

	"github.com/gin-gonic/gin"
)

type GameController struct {
	GameService *service.GameService
}

func NewGameController() *GameController {
	return &GameController{
		GameService: service.NewGameService(),
	}
}

func (ctrl *GameController) GetGameList(c *gin.Context) {
	list := ctrl.GameService.GetGameList()
	response.Success(c, list)
}

func (ctrl *GameController) CreateGame(c *gin.Context) {
	var param dto.GameRoomParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	userUuid := c.GetHeader("userUuid")
	if userUuid == "" {
		response.Failed(c, 401, "User not identified")
		return
	}

	info, err := ctrl.GameService.CreateGame(userUuid, param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, info)
}

func (ctrl *GameController) JoinGame(c *gin.Context) {
	var param dto.GameRoomParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	userUuid := c.GetHeader("userUuid")
	if userUuid == "" {
		response.Failed(c, 401, "User not identified")
		return
	}

	member, err := ctrl.GameService.JoinGame(userUuid, param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, member)
}

func (ctrl *GameController) StartGame(c *gin.Context) {
	var param dto.GameRoomParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	userUuid := c.GetHeader("userUuid")
	if userUuid == "" {
		response.Failed(c, 401, "User not identified")
		return
	}

	err := ctrl.GameService.StartGame(userUuid, param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *GameController) EndGame(c *gin.Context) {
	var param dto.GameRoomParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Failed(c, 400, "Invalid parameters")
		return
	}

	userUuid := c.GetHeader("userUuid")
	if userUuid == "" {
		response.Failed(c, 401, "User not identified")
		return
	}

	err := ctrl.GameService.EndGame(userUuid, param)
	if err != nil {
		response.Failed(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}
