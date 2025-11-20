package dto

type SingParam struct {
	RoomUuid string `json:"roomUuid"`
	OrderId  uint64 `json:"orderId"`
	ChorusId string `json:"chorusId"`
	Ext      string `json:"ext"`
}

type SingInfoParam struct {
	RoomUuid string `json:"roomUuid"`
}

type SingActionParam struct {
	RoomUuid string `json:"roomUuid"`
	Action   int    `json:"action"` // 1: Pause, 2: Resume, 3: End
}

type ChorusInviteParam struct {
	RoomUuid string `json:"roomUuid"`
	OrderId  uint64 `json:"orderId"`
	UserUuid string `json:"userUuid"` // Invitee
}

type JoinChorusParam struct {
	RoomUuid string `json:"roomUuid"`
	ChorusId string `json:"chorusId"`
}

type CancelChorusParam struct {
	RoomUuid string `json:"roomUuid"`
	ChorusId string `json:"chorusId"`
}

type FinishChorusReadyParam struct {
	RoomUuid string `json:"roomUuid"`
	ChorusId string `json:"chorusId"`
}

type AbandonSingParam struct {
	RoomUuid string `json:"roomUuid"`
	OrderId  uint64 `json:"orderId"`
}

type OrderSongParam struct {
	LiveRecordId uint64 `json:"liveRecordId"`
	SongId       string `json:"songId"`
	SongName     string `json:"songName"`
	SongCover    string `json:"songCover"`
	Singer       string `json:"singer"`
	SingerCover  string `json:"singerCover"`
	SongTime     int64  `json:"songTime"`
	Channel      int    `json:"channel"`
}

type SwitchSongParam struct {
	LiveRecordId   uint64 `json:"liveRecordId"`
	CurrentOrderId uint64 `json:"currentOrderId"`
	Attachment     string `json:"attachment"`
}

type OrderParam struct {
	LiveRecordId uint64 `json:"liveRecordId"`
	OrderId      uint64 `json:"orderId"`
}

type CleanOrderSongParam struct {
	LiveRecordId uint64 `json:"liveRecordId"`
}

type ChorusControlResultDto struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type SingDetailInfoDto struct {
	Song *OrderSongDto `json:"song"`
}

type OrderSongDto struct {
	Id            uint64 `json:"id"`
	LiveRecordId  uint64 `json:"liveRecordId"`
	RoomArchiveId string `json:"roomArchiveId"`
	RoomUuid      string `json:"roomUuid"`
	UserUuid      string `json:"userUuid"`
	UserName      string `json:"userName"`
	UserIcon      string `json:"userIcon"`
	SongId        string `json:"songId"`
	SongName      string `json:"songName"`
	SongCover     string `json:"songCover"`
	Singer        string `json:"singer"`
	SingerCover   string `json:"singerCover"`
	SongTime      int64  `json:"songTime"`
	Channel       int    `json:"channel"`
	Status        int    `json:"status"`
	SetTopTime    int64  `json:"setTopTime"`
}
