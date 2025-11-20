package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Spring   SpringConfig   `mapstructure:"spring"`
	Yunxin   YunxinConfig   `mapstructure:"yunxin"`
	Business BusinessConfig `mapstructure:"business"`
	Server   ServerConfig   `mapstructure:"server"`
}

type SpringConfig struct {
	Datasource DatasourceConfig `mapstructure:"datasource"`
	Redis      RedisConfig      `mapstructure:"redis"`
}

type DatasourceConfig struct {
	Driver   string       `mapstructure:"driver-class-name"`
	Hikari   HikariConfig `mapstructure:"hikari"`
	Url      string       `mapstructure:"url"`
	Username string       `mapstructure:"username"`
	Password string       `mapstructure:"password"`
}

type HikariConfig struct {
	MaximumPoolSize int `mapstructure:"maximum-pool-size"`
}

type RedisConfig struct {
	Database int    `mapstructure:"database"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type YunxinConfig struct {
	Origin OriginConfig `mapstructure:"origin"`
}

type OriginConfig struct {
	AppKey            string `mapstructure:"appKey"`
	AppSecret         string `mapstructure:"appSecret"`
	NimHost           string `mapstructure:"nimHost"`
	NeLiveHost        string `mapstructure:"neLiveHost"`
	NeRoomHost        string `mapstructure:"neRoomHost"`
	SecurityAuditHost string `mapstructure:"securityAuditHost"`
	RtcHost           string `mapstructure:"rtcHost"`
}

type BusinessConfig struct {
	YunxinAssistAccid      string     `mapstructure:"yunxinAssistAccid"`
	SystemAccid            string     `mapstructure:"systemAccid"`
	OneVOneRtcRoomLiveTime int        `mapstructure:"1v1RtcRoomLiveTime"`
	VoiceRoomConfigId      int        `mapstructure:"voiceRoomConfigId"`
	GameRoomConfigId       int        `mapstructure:"gameRoomConfigId"`
	ListenTogetherConfigId int        `mapstructure:"listenTogetherConfigId"`
	KtvConfigId            int        `mapstructure:"ktvConfigId"`
	PkConfigId             int        `mapstructure:"pkConfigId"`
	RoomOrderSongLimit     int        `mapstructure:"roomOrderSongLimit"`
	UserOrderSongLimit     int        `mapstructure:"userOrderSongLimit"`
	Game                   GameConfig `mapstructure:"game"`
}

type GameConfig struct {
	SudUrl      string `mapstructure:"sudUrl"`
	AppId       string `mapstructure:"appId"`
	AppKey      string `mapstructure:"appKey"`
	AppSecret   string `mapstructure:"appSecret"`
	OnlineGames string `mapstructure:"onlineGames"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("application-local")
	viper.SetConfigType("yaml")
	// Search paths:
	// 1. Original Java resources (for easy dev)
	// 2. Local config folder
	// 3. Current directory
	viper.AddConfigPath("../../nemo-controller/src/main/resources")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}
