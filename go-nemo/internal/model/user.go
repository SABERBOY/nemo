package model

import "time"

type User struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Mobile     string    `gorm:"type:varchar(15);default:''" json:"mobile"`
	UserUuid   string    `gorm:"type:varchar(64);not null;unique" json:"userUuid"`
	UserToken  string    `gorm:"type:varchar(64);not null" json:"userToken"`
	ImToken    string    `gorm:"type:varchar(64);default:''" json:"imToken"`
	UserName   string    `gorm:"type:varchar(64);default:''" json:"userName"`
	Icon       string    `gorm:"type:varchar(1024);default:''" json:"icon"`
	Age        int       `gorm:"default:-1" json:"age"`
	Sex        int       `gorm:"default:0" json:"sex"`   // 0 unknown, 1 male, 2 female
	State      int       `gorm:"default:1" json:"state"` // 1 available, 2 disabled
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type UserDevice struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserUuid   string    `gorm:"type:varchar(64);not null;uniqueIndex:user_device" json:"userUuid"`
	DeviceId   string    `gorm:"type:varchar(64);not null;uniqueIndex:user_device;index:device" json:"deviceId"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type Gift struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	GiftName   string    `gorm:"type:varchar(64);not null" json:"giftName"`
	GiftDesc   string    `gorm:"type:varchar(64);not null" json:"giftDesc"`
	CloudCoin  int64     `json:"cloudCoin"`
	Status     int       `gorm:"default:1" json:"status"` // -1 invalid, 1 valid
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type UserReward struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserUuid   string    `gorm:"type:varchar(64);not null" json:"userUuid"`
	GiftId     uint64    `gorm:"not null" json:"giftId"`
	GiftCount  int       `gorm:"default:1" json:"giftCount"`
	CloudCoin  int64     `json:"cloudCoin"`
	Target     string    `gorm:"type:varchar(64);not null" json:"target"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}
