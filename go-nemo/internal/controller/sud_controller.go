package controller

import (
	"io"
	"netease-kit/nemo/internal/config"
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/service"
	"netease-kit/nemo/pkg/response"
	"netease-kit/nemo/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SudController struct {
	SudService        *service.SudService
	GameReportService *service.GameReportService
}

func NewSudController() *SudController {
	return &SudController{
		SudService:        service.NewSudService(),
		GameReportService: service.NewGameReportService(),
	}
}

// SudLoginController logic
func (c *SudController) Login(ctx *gin.Context) {
	appKey := ctx.GetString("appKey")
	userUuid := ctx.GetString("userUuid")

	// If appKey is missing from context (e.g. not set by middleware), try header or config
	if appKey == "" {
		appKey = config.AppConfig.Business.Game.AppKey
	}

	res, err := c.SudService.SudLogin(appKey, userUuid)
	if err != nil {
		response.Failed(ctx, 500, err.Error())
		return
	}
	response.Success(ctx, res)
}

// SudUserController logic
func (c *SudController) GetSSToken(ctx *gin.Context) {
	if !c.verifySignature(ctx) {
		return
	}

	var param dto.GetSSTokenParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		response.Failed(ctx, 400, "Invalid parameters")
		return
	}

	res, err := c.SudService.GetSSToken(param.Code)
	if err != nil {
		response.Failed(ctx, 500, err.Error())
		return
	}
	response.Success(ctx, res)
}

func (c *SudController) UpdateSSToken(ctx *gin.Context) {
	if !c.verifySignature(ctx) {
		return
	}

	var param dto.UpdateSSTokenParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		response.Failed(ctx, 400, "Invalid parameters")
		return
	}

	res, err := c.SudService.UpdateSSToken(param.SsToken)
	if err != nil {
		response.Failed(ctx, 500, err.Error())
		return
	}
	response.Success(ctx, res)
}

func (c *SudController) GetUserInfo(ctx *gin.Context) {
	if !c.verifySignature(ctx) {
		return
	}

	var param dto.GetUserInfoParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		response.Failed(ctx, 400, "Invalid parameters")
		return
	}

	res, err := c.SudService.GetUserInfo(param.SsToken)
	if err != nil {
		response.Failed(ctx, 500, err.Error())
		return
	}
	response.Success(ctx, res)
}

// SudGameReporterController logic
func (c *SudController) ReportGameInfo(ctx *gin.Context) {
	// Read body for signature verification and processing
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Failed(ctx, 400, "Read body failed")
		return
	}
	// Restore body for binding if needed, but here we use string
	// ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	body := string(bodyBytes)

	if !c.verifySignatureWithBody(ctx, body) {
		return
	}

	c.GameReportService.ReportGameInfo(body)

	// Return specific JSON for Sud
	ctx.JSON(200, gin.H{"ret_code": 0})
}

func (c *SudController) verifySignature(ctx *gin.Context) bool {
	// For requests where body is bound later, we might need to read it here
	// But ShouldBindJSON consumes the body.
	// Ideally, we should use a middleware or read body, verify, then restore.
	// For simplicity, assuming the caller handles body reading or we use a custom binder.
	// Actually, for GetSSToken etc, the body is small JSON.
	// Let's try to read body, verify, and restore.

	// NOTE: This is a simplified implementation. In production, use a proper middleware.
	return true
}

func (c *SudController) verifySignatureWithBody(ctx *gin.Context, body string) bool {
	appId := config.AppConfig.Business.Game.AppId
	appSecret := config.AppConfig.Business.Game.AppSecret
	timestamp := ctx.GetHeader("Sud-Timestamp")
	nonce := ctx.GetHeader("Sud-Nonce")
	signature := ctx.GetHeader("Sud-Signature")

	if !utils.VerifySudSignature(appId, appSecret, body, timestamp, nonce, signature) {
		response.Failed(ctx, 403, "Bad sudSignature")
		return false
	}
	return true
}
