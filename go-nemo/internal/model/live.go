package model

import "time"

type LiveRecord struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomArchiveId string    `gorm:"type:varchar(64);index:idx_room_archive_id" json:"roomArchiveId"`
	RoomUuid      string    `gorm:"type:varchar(64);index:idx_room_uuid" json:"roomUuid"`
	RoomName      string    `gorm:"type:varchar(64)" json:"roomName"`
	UserUuid      string    `gorm:"type:varchar(64);index:idx_user_uuid_status" json:"userUuid"`
	LiveTopic     string    `gorm:"type:varchar(256)" json:"liveTopic"`
	Cover         string    `gorm:"type:varchar(256)" json:"cover"`
	Status        int       `gorm:"default:1;index:idx_user_uuid_status" json:"status"` // 1: valid, -1: invalid
	Live          int       `gorm:"default:1" json:"live"`                              // -1: ended, 0: not started, 1: living
	LiveType      int       `gorm:"default:1" json:"liveType"`                          // 1: Interactive Live, 2: Voice Room, 3: KTV
	LiveConfig    string    `gorm:"type:varchar(1024)" json:"liveConfig"`
	SingMode      int       `gorm:"default:0" json:"singMode"` // 0: Smart Chorus, 1: Serial, 2: NTP Realtime, 3: Solo
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type LiveReward struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	LiveRecordId  uint64    `gorm:"not null;index:idx_live_record_id_user_uuid;index:idx_live_record_id_target" json:"liveRecordId"`
	RoomArchiveId string    `gorm:"type:varchar(64)" json:"roomArchiveId"`
	RoomUuid      string    `gorm:"type:varchar(64);not null;default:''" json:"roomUuid"`
	UserUuid      string    `gorm:"type:varchar(64);not null;index:idx_live_record_id_user_uuid" json:"userUuid"`
	GiftId        uint64    `json:"giftId"`
	CloudCoin     int64     `json:"cloudCoin"`
	GiftCount     int       `gorm:"default:1" json:"giftCount"`
	Target        string    `gorm:"type:varchar(64);not null;index:idx_live_record_id_target" json:"target"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type OrderSong struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	LiveRecordId  uint64    `gorm:"index:idx_live_record_id_useruuid" json:"liveRecordId"`
	RoomArchiveId string    `gorm:"type:varchar(64)" json:"roomArchiveId"`
	RoomUuid      string    `gorm:"type:varchar(64);not null;default:''" json:"roomUuid"`
	UserUuid      string    `gorm:"type:varchar(64);not null;default:'';index:idx_live_record_id_useruuid" json:"userUuid"`
	SongId        string    `gorm:"type:varchar(64);not null;default:''" json:"songId"`
	SongName      string    `gorm:"type:varchar(64)" json:"songName"`
	SongCover     string    `gorm:"type:varchar(512)" json:"songCover"`
	Singer        string    `gorm:"type:varchar(64)" json:"singer"`
	SingerCover   string    `gorm:"type:varchar(512)" json:"singerCover"`
	SongTime      int64     `gorm:"default:0" json:"songTime"`
	Channel       int       `gorm:"default:1" json:"channel"` // 1: Cloud Music, 2: Migu
	Status        int       `gorm:"default:0" json:"status"`  // -2: Sung, -1: Deleted, 0: Waiting, 1: Singing
	SetTopTime    int64     `gorm:"default:0" json:"setTopTime"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type ChorusRecord struct {
	Id                   uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ChorusId             string    `gorm:"type:varchar(64);not null;default:'';unique" json:"chorusId"`
	RoomUuid             string    `gorm:"type:varchar(64);not null;default:''" json:"roomUuid"`
	RoomName             string    `gorm:"type:varchar(64);not null;default:''" json:"roomName"`
	LiveRecordId         uint64    `gorm:"not null" json:"liveRecordId"`
	LeaderUuid           string    `gorm:"type:varchar(64);not null;default:''" json:"leaderUuid"`
	AssistantUuid        string    `gorm:"type:varchar(64);not null;default:''" json:"assistantUuid"`
	Status               int       `gorm:"default:1" json:"status"` // 1: Valid, -1: Invalid
	State                int       `gorm:"default:0" json:"state"`  // 0: Inviting, 1: Singing, 2: Rejected, 3: Cancelled, 4: Ended, 5: Agreed
	LeaderDeviceParam    string    `gorm:"type:varchar(128);not null;default:''" json:"leaderDeviceParam"`
	AssistantDeviceParam string    `gorm:"type:varchar(128);not null;default:''" json:"assistantDeviceParam"`
	OrderId              uint64    `gorm:"default:0" json:"orderId"`
	CreateTime           time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime           time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}
