package controller

import (
	"netease-kit/nemo/internal/repo"
	"netease-kit/nemo/pkg/response"

	"github.com/gin-gonic/gin"
)

type RedisQueueController struct {
	queueRepo *repo.QueueRepository
}

func NewRedisQueueController() *RedisQueueController {
	return &RedisQueueController{
		queueRepo: repo.NewQueueRepository(),
	}
}

type SendQueueRequest struct {
	Topic   string `form:"topic" binding:"required"`
	Message string `form:"message" binding:"required"`
	Delay   int64  `form:"delay" binding:"required"`
}

type CancelQueueRequest struct {
	Topic   string `form:"topic" binding:"required"`
	Message string `form:"message" binding:"required"`
}

func (c *RedisQueueController) Send(ctx *gin.Context) {
	var req SendQueueRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Failed(ctx, 400, "Invalid parameters: "+err.Error())
		return
	}

	err := c.queueRepo.SendDelayMessage(req.Topic, req.Message, req.Delay)
	if err != nil {
		response.Failed(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, nil)
}

func (c *RedisQueueController) Cancel(ctx *gin.Context) {
	var req CancelQueueRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Failed(ctx, 400, "Invalid parameters: "+err.Error())
		return
	}

	err := c.queueRepo.CancelDelayMessage(req.Topic, req.Message)
	if err != nil {
		response.Failed(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, nil)
}
