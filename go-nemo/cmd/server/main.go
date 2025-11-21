package main

import (
	"fmt"
	"netease-kit/nemo/internal/config"
	"netease-kit/nemo/internal/db"
	"time"

	"netease-kit/nemo/internal/controller"
	"netease-kit/nemo/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	// 1. Init Config
	config.InitConfig()

	// 2. Init DB
	db.InitDB()

	// 3. Init Router
	r := gin.Default()

	// 4. Global Middleware
	r.Use(ContextMiddleware())

	// 5. Routes
	initController := controller.NewNemoInitController()
	entLiveController := controller.NewEntLiveController()
	socialChatController := controller.NewSocialChatController()
	gameController := controller.NewGameController()
	ktvController := controller.NewKtvController()
	songController := controller.NewSongController()
	musicPlayerController := controller.NewMusicPlayerController()
	sudController := controller.NewSudController()

	// Notify Controller Dependencies
	entLiveService := service.NewEntLiveService()
	musicService := service.NewMusicPlayService()
	ktvService := service.NewKtvService()
	imEventService := service.NewImEventService()
	notifyService := service.NewNotifyService(entLiveService, musicService, ktvService)
	notifyController := controller.NewNotifyController(notifyService, imEventService)
	redisQueueController := controller.NewRedisQueueController()

	nemoGroup := r.Group("/nemo/app")
	{
		nemoGroup.POST("/initAppAndUser", initController.InitAppAndUser)
	}

	entLiveGroup := r.Group("/nemo/entertainmentLive")
	{
		entLiveGroup.POST("/createLive", entLiveController.CreateLive)
		entLiveGroup.POST("/closeLive", entLiveController.CloseLive)
		entLiveGroup.GET("/list", entLiveController.GetLiveList)
		entLiveGroup.GET("/info", entLiveController.GetLiveInfo)

		// KTV Routes
		ktvGroup := entLiveGroup.Group("/ktv")
		{
			ktvGroup.POST("/sing/start", ktvController.SingStart)
			ktvGroup.POST("/sing/info", ktvController.GetSingInfo)
			ktvGroup.POST("/sing/action", ktvController.SingControl)
			ktvGroup.POST("/sing/chorus/invite", ktvController.ChorusInvite)
			ktvGroup.POST("/sing/chorus/join", ktvController.ChorusJoin)
			ktvGroup.POST("/sing/chorus/cancel", ktvController.ChorusCancel)
			ktvGroup.POST("/sing/chorus/ready", ktvController.ChorusReady)
			ktvGroup.POST("/sing/abandon", ktvController.AbandonSing)
		}

		// Song Routes
		songGroup := entLiveGroup.Group("/live/song")
		{
			songGroup.POST("/orderSong", songController.OrderSong)
			songGroup.POST("/switchSong", songController.SwitchSong)
			songGroup.GET("/getOrderSongs", songController.GetOrderSongs)
			songGroup.POST("/cancelOrderSong", songController.CancelOrderSong)
			songGroup.POST("/songSetTop", songController.SongSetTop)
			songGroup.POST("/cleanUserOrderSongs", songController.CleanUserOrderSongs)
			songGroup.POST("/getMusicToken", songController.GetMusicToken)
		}

		// Music Player Routes
		musicGroup := entLiveGroup.Group("/music")
		{
			musicGroup.GET("/info", musicPlayerController.GetPlayMusicInfo)
			musicGroup.POST("/action", musicPlayerController.MusicAction)
			musicGroup.POST("/ready", musicPlayerController.MusicReady)
		}

		// Notify Routes
		notifyGroup := entLiveGroup.Group("/nim")
		{
			notifyGroup.POST("/notify", notifyController.Notify)
			notifyGroup.POST("/im-event-notify", notifyController.ImEventNotify)
		}
	}

	socialChatGroup := r.Group("/nemo/socialChat")
	{
		socialChatGroup.POST("/user/reporter", socialChatController.Reporter)
		socialChatGroup.GET("/user/getOnLineUser", socialChatController.GetOnLineUser)
		socialChatGroup.POST("/user/reward", socialChatController.UserReward)
		socialChatGroup.POST("/user/login", socialChatController.Login)
		socialChatGroup.POST("/user/getUserState", socialChatController.GetUserState)
		socialChatGroup.POST("/user/getUserInfo", socialChatController.GetUserInfo)
	}

	gameGroup := r.Group("/nemo/game")
	{
		gameGroup.GET("/list", gameController.GetGameList)
		gameGroup.POST("/create", gameController.CreateGame)
		gameGroup.POST("/join", gameController.JoinGame)
		gameGroup.POST("/start", gameController.StartGame)
		gameGroup.POST("/end", gameController.EndGame)
		gameGroup.POST("/exit", gameController.ExitGame)
		gameGroup.GET("/members", gameController.GetGameMembers)
		gameGroup.GET("/gameInfo", gameController.GetGameInfo)
		gameGroup.GET("/status-reporter", gameController.StatusReporter)

		sudGroup := gameGroup.Group("/sud")
		{
			sudGroup.POST("/login", sudController.Login)

			userGroup := sudGroup.Group("/user")
			{
				userGroup.POST("/get_sstoken", sudController.GetSSToken)
				userGroup.POST("/update_sstoken", sudController.UpdateSSToken)
				userGroup.POST("/get_user_info", sudController.GetUserInfo)
				userGroup.POST("/report_game_info", sudController.ReportGameInfo)
			}
		}
	}

	redisQueueGroup := r.Group("/nemo/redis-queue")
	{
		redisQueueGroup.POST("/send", redisQueueController.Send)
		redisQueueGroup.POST("/cancel", redisQueueController.Cancel)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 6. Run
	port := config.AppConfig.Server.Port
	if port == 0 {
		port = 9981 // Default fallback
	}
	r.Run(fmt.Sprintf(":%d", port))
}

func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start Time
		c.Set("startTime", time.Now())

		// Request ID
		requestId := c.GetHeader("requestId")
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Set("requestId", requestId)

		// AppKey & Secret
		appKey := c.GetHeader("AppKey")
		appSecret := c.GetHeader("AppSecret")
		c.Set("appKey", appKey)
		c.Set("appSecret", appSecret)

		c.Next()
	}
}
