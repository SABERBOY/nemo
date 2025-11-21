package model

import "time"

type GameRecord struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	LiveRecordId  uint64    `gorm:"index:idx_live_record_id" json:"liveRecordId"`
	RoomUuid      string    `gorm:"type:varchar(64);index:idx_room_uuid" json:"roomUuid"`
	GameCreator   string    `gorm:"type:varchar(64)" json:"gameCreator"`
	GameStatus    int       `gorm:"default:0" json:"gameStatus"` // 0: Waiting, 1: Playing, 2: Ended
	GameId        uint64    `json:"gameId"`
	GameName      string    `gorm:"type:varchar(64)" json:"gameName"`
	GameDesc      string    `gorm:"type:varchar(64)" json:"gameDesc"`
	Thumbnail     string    `gorm:"type:varchar(256)" json:"thumbnail"`
	Rule          string    `gorm:"type:varchar(256)" json:"rule"`
	RoomArchiveId string    `gorm:"type:varchar(64)" json:"roomArchiveId"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type GameMember struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	LiveRecordId  uint64    `gorm:"index:idx_live_record_id" json:"liveRecordId"`
	RoomUuid      string    `gorm:"type:varchar(64);index:idx_gameId_room_uuid" json:"roomUuid"`
	Status        int       `gorm:"default:0" json:"status"` // 0: Preparing, 1: Playing, 2: Left
	JoinTime      int64     `gorm:"default:0" json:"joinTime"`
	ExitTime      int64     `gorm:"default:0" json:"exitTime"`
	GameId        uint64    `gorm:"index:idx_gameId_room_uuid" json:"gameId"`
	UserUuid      string    `gorm:"type:varchar(64)" json:"userUuid"`
	UserName      string    `gorm:"type:varchar(64)" json:"userName"`
	RoomArchiveId string    `gorm:"type:varchar(64)" json:"roomArchiveId"`
	GameRecordId  uint64    `json:"gameRecordId"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type GameReport struct {
	Id           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	GameRecordId uint64    `json:"gameRecordId"`
	ReportMsg    string    `gorm:"type:text" json:"reportMsg"`
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`
}
