package dto

type OnLineUserDto struct {
	UserUuid        string `json:"userUuid"`
	UserName        string `json:"userName"`
	Icon            string `json:"icon"`
	Mobile          string `json:"mobile"`
	FirstReportTime int64  `json:"firstReportTime"`
	LastReportTime  int64  `json:"lastReportTime"`
}

type UserRewardDto struct {
	UserUuid  string `json:"userUuid"`
	Target    string `json:"target"`
	GiftId    uint64 `json:"giftId"`
	GiftCount int    `json:"giftCount"`
}

type RtcRoomNotifyParam struct {
	ChannelId   int64  `json:"channelId"`
	ChannelName string `json:"channelName"`
	Status      int    `json:"status"` // 1: Start, 2: End
	Timestamp   int64  `json:"timestamp"`
}

type RtcRoomUserNotifyParam struct {
	ChannelId   int64  `json:"channelId"`
	ChannelName string `json:"channelName"`
	Uid         int64  `json:"uid"`
	Status      int    `json:"status"` // 1: Join, 2: Leave
	Timestamp   int64  `json:"timestamp"`
}

type RtcRoomInfoDto struct {
	ChannelId   int64  `json:"channelId"`
	ChannelName string `json:"channelName"`
	Status      int    `json:"status"`
	AppKey      string `json:"appKey"`
}

type RtcRoomUserInfoDto struct {
	ChannelId   int64  `json:"channelId"`
	ChannelName string `json:"channelName"`
	Uid         int64  `json:"uid"`
	UserUuid    string `json:"userUuid"`
	UserName    string `json:"userName"`
	Icon        string `json:"icon"`
	AppKey      string `json:"appKey"`
}

type EventDto struct {
	Data interface{} `json:"data"`
	Type int         `json:"type"`
}

type RewardMessage struct {
	SenderUserUuid string `json:"senderUserUuid"`
	TargetUserUuid string `json:"targetUserUuid"`
	GiftId         uint64 `json:"giftId"`
	GiftCount      int    `json:"giftCount"`
	CloudCoin      int64  `json:"cloudCoin"`
}
